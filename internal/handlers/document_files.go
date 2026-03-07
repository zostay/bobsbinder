package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DocumentFileHandler struct {
	DB        *sql.DB
	Logger    *zap.Logger
	UploadDir string
}

func (h *DocumentFileHandler) verifyDocOwnership(docID, userID int64) error {
	var count int
	err := h.DB.QueryRow(`
		SELECT COUNT(*) FROM documents d
		JOIN parties p ON d.party_id = p.id
		WHERE d.id = ? AND p.user_id = ? AND d.deleted_at IS NULL
	`, docID, userID).Scan(&count)
	if err != nil || count == 0 {
		return fmt.Errorf("document not found")
	}
	return nil
}

func (h *DocumentFileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	docID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid document id"}`, http.StatusBadRequest)
		return
	}

	if err := h.verifyDocOwnership(docID, userID); err != nil {
		http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error":"file too large"}`, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"no file provided"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	dir := filepath.Join(h.UploadDir, strconv.FormatInt(userID, 10), strconv.FormatInt(docID, 10))
	if err := os.MkdirAll(dir, 0750); err != nil {
		h.Logger.Error("failed to create upload directory", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	storedName := uuid.New().String() + "_" + filepath.Base(header.Filename)
	filePath := filepath.Join(dir, storedName)

	dst, err := os.Create(filePath)
	if err != nil {
		h.Logger.Error("failed to create file", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		h.Logger.Error("failed to write file", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	result, err := h.DB.Exec(
		"INSERT INTO document_files (document_id, filename, content_type, file_path, file_size) VALUES (?, ?, ?, ?, ?)",
		docID, header.Filename, contentType, filePath, header.Size,
	)
	if err != nil {
		h.Logger.Error("failed to insert document file", zap.Error(err))
		os.Remove(filePath)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	fileID, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":           fileID,
		"document_id":  docID,
		"filename":     header.Filename,
		"content_type": contentType,
		"file_size":    header.Size,
	})
}

func (h *DocumentFileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	docID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid document id"}`, http.StatusBadRequest)
		return
	}

	if err := h.verifyDocOwnership(docID, userID); err != nil {
		http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
		return
	}

	rows, err := h.DB.Query(
		"SELECT id, document_id, filename, content_type, file_size, created_at FROM document_files WHERE document_id = ? ORDER BY created_at DESC",
		docID,
	)
	if err != nil {
		h.Logger.Error("failed to list document files", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var files []map[string]any
	for rows.Next() {
		var id, documentID, fileSize int64
		var filename, contentType, createdAt string
		if err := rows.Scan(&id, &documentID, &filename, &contentType, &fileSize, &createdAt); err != nil {
			h.Logger.Error("failed to scan document file", zap.Error(err))
			continue
		}
		files = append(files, map[string]any{
			"id":           id,
			"document_id":  documentID,
			"filename":     filename,
			"content_type": contentType,
			"file_size":    fileSize,
			"created_at":   createdAt,
		})
	}

	if files == nil {
		files = []map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func (h *DocumentFileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	docID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid document id"}`, http.StatusBadRequest)
		return
	}

	fileID, err := strconv.ParseInt(chi.URLParam(r, "fileId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid file id"}`, http.StatusBadRequest)
		return
	}

	if err := h.verifyDocOwnership(docID, userID); err != nil {
		http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
		return
	}

	var filename, filePath string
	err = h.DB.QueryRow(
		"SELECT filename, file_path FROM document_files WHERE id = ? AND document_id = ?",
		fileID, docID,
	).Scan(&filename, &filePath)
	if err != nil {
		http.Error(w, `{"error":"file not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	http.ServeFile(w, r, filePath)
}

func (h *DocumentFileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	docID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid document id"}`, http.StatusBadRequest)
		return
	}

	fileID, err := strconv.ParseInt(chi.URLParam(r, "fileId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid file id"}`, http.StatusBadRequest)
		return
	}

	if err := h.verifyDocOwnership(docID, userID); err != nil {
		http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
		return
	}

	var filePath string
	err = h.DB.QueryRow(
		"SELECT file_path FROM document_files WHERE id = ? AND document_id = ?",
		fileID, docID,
	).Scan(&filePath)
	if err != nil {
		http.Error(w, `{"error":"file not found"}`, http.StatusNotFound)
		return
	}

	_, err = h.DB.Exec("DELETE FROM document_files WHERE id = ?", fileID)
	if err != nil {
		h.Logger.Error("failed to delete document file record", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	os.Remove(filePath)

	w.WriteHeader(http.StatusNoContent)
}
