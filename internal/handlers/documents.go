package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type DocumentHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type documentRequest struct {
	PartyID     int64  `json:"party_id"`
	CategoryID  int64  `json:"category_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Status      string `json:"status"`
	DocType     string `json:"doc_type"`
	SecureNotes string `json:"secure_notes"`
	LocationID  *int64 `json:"location_id"`
}

func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(`
		SELECT d.id, d.party_id, d.category_id, d.title, d.content, d.status, d.doc_type, d.location_id, d.secure_notes, d.created_at, d.updated_at
		FROM documents d
		JOIN parties p ON d.party_id = p.id
		WHERE p.user_id = ? AND d.deleted_at IS NULL
		ORDER BY d.created_at DESC
	`, userID)
	if err != nil {
		h.Logger.Error("failed to list documents", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var documents []map[string]any
	for rows.Next() {
		var id, partyID, categoryID int64
		var title, content, status, secureNotes string
		var docType sql.NullString
		var locationID sql.NullInt64
		var createdAt, updatedAt string
		if err := rows.Scan(&id, &partyID, &categoryID, &title, &content, &status, &docType, &locationID, &secureNotes, &createdAt, &updatedAt); err != nil {
			h.Logger.Error("failed to scan document", zap.Error(err))
			continue
		}
		doc := map[string]any{
			"id":           id,
			"party_id":     partyID,
			"category_id":  categoryID,
			"title":        title,
			"content":      content,
			"status":       status,
			"doc_type":     docType.String,
			"secure_notes": secureNotes,
			"created_at":   createdAt,
			"updated_at":   updatedAt,
		}
		if locationID.Valid {
			doc["location_id"] = locationID.Int64
		} else {
			doc["location_id"] = nil
		}
		documents = append(documents, doc)
	}

	if documents == nil {
		documents = []map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

func (h *DocumentHandler) Get(w http.ResponseWriter, r *http.Request) {
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

	var id, partyID, categoryID int64
	var title, content, status, secureNotes string
	var docType sql.NullString
	var locationID sql.NullInt64
	var createdAt, updatedAt string
	err = h.DB.QueryRow(`
		SELECT d.id, d.party_id, d.category_id, d.title, d.content, d.status, d.doc_type, d.location_id, d.secure_notes, d.created_at, d.updated_at
		FROM documents d
		JOIN parties p ON d.party_id = p.id
		WHERE d.id = ? AND p.user_id = ? AND d.deleted_at IS NULL
	`, docID, userID).Scan(&id, &partyID, &categoryID, &title, &content, &status, &docType, &locationID, &secureNotes, &createdAt, &updatedAt)
	if err != nil {
		http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
		return
	}

	doc := map[string]any{
		"id":           id,
		"party_id":     partyID,
		"category_id":  categoryID,
		"title":        title,
		"content":      content,
		"status":       status,
		"doc_type":     docType.String,
		"secure_notes": secureNotes,
		"created_at":   createdAt,
		"updated_at":   updatedAt,
	}
	if locationID.Valid {
		doc["location_id"] = locationID.Int64
	} else {
		doc["location_id"] = nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (h *DocumentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req documentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Auto-resolve party_id to user's "self" party when omitted
	if req.PartyID == 0 {
		err := h.DB.QueryRow(
			"SELECT id FROM parties WHERE user_id = ? AND relationship = 'self' LIMIT 1",
			userID,
		).Scan(&req.PartyID)
		if err != nil {
			http.Error(w, `{"error":"no self party found"}`, http.StatusBadRequest)
			return
		}
	}

	// Verify party belongs to user
	var count int
	err := h.DB.QueryRow("SELECT COUNT(*) FROM parties WHERE id = ? AND user_id = ?", req.PartyID, userID).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, `{"error":"party not found"}`, http.StatusNotFound)
		return
	}

	if req.Status == "" {
		req.Status = "draft"
	}
	if req.DocType == "" {
		req.DocType = "reference"
	}

	result, err := h.DB.Exec(
		"INSERT INTO documents (party_id, category_id, title, content, status, doc_type, location_id, secure_notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		req.PartyID, req.CategoryID, req.Title, req.Content, req.Status, req.DocType, req.LocationID, req.SecureNotes,
	)
	if err != nil {
		h.Logger.Error("failed to create document", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	docID, _ := result.LastInsertId()

	resp := map[string]any{
		"id":           docID,
		"party_id":     req.PartyID,
		"category_id":  req.CategoryID,
		"title":        req.Title,
		"content":      req.Content,
		"status":       req.Status,
		"doc_type":     req.DocType,
		"secure_notes": req.SecureNotes,
	}
	if req.LocationID != nil {
		resp["location_id"] = *req.LocationID
	} else {
		resp["location_id"] = nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *DocumentHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req documentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE documents d
		JOIN parties p ON d.party_id = p.id
		SET d.title = ?, d.content = ?, d.status = ?, d.doc_type = ?, d.location_id = ?, d.category_id = ?, d.secure_notes = ?
		WHERE d.id = ? AND p.user_id = ? AND d.deleted_at IS NULL
	`, req.Title, req.Content, req.Status, req.DocType, req.LocationID, req.CategoryID, req.SecureNotes, docID, userID)
	if err != nil {
		h.Logger.Error("failed to update document", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *DocumentHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.DB.Exec(`
		UPDATE documents d
		JOIN parties p ON d.party_id = p.id
		SET d.deleted_at = NOW()
		WHERE d.id = ? AND p.user_id = ? AND d.deleted_at IS NULL
	`, docID, userID)
	if err != nil {
		h.Logger.Error("failed to delete document", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DocumentHandler) ListTrash(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(`
		SELECT d.id, d.party_id, d.category_id, d.title, d.content, d.status, d.doc_type, d.location_id, d.secure_notes, d.created_at, d.updated_at, d.deleted_at
		FROM documents d
		JOIN parties p ON d.party_id = p.id
		WHERE p.user_id = ? AND d.deleted_at IS NOT NULL
		ORDER BY d.deleted_at DESC
	`, userID)
	if err != nil {
		h.Logger.Error("failed to list trashed documents", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var documents []map[string]any
	for rows.Next() {
		var id, partyID, categoryID int64
		var title, content, status, secureNotes string
		var docType sql.NullString
		var locationID sql.NullInt64
		var createdAt, updatedAt, deletedAt string
		if err := rows.Scan(&id, &partyID, &categoryID, &title, &content, &status, &docType, &locationID, &secureNotes, &createdAt, &updatedAt, &deletedAt); err != nil {
			h.Logger.Error("failed to scan trashed document", zap.Error(err))
			continue
		}
		doc := map[string]any{
			"id":           id,
			"party_id":     partyID,
			"category_id":  categoryID,
			"title":        title,
			"content":      content,
			"status":       status,
			"doc_type":     docType.String,
			"secure_notes": secureNotes,
			"created_at":   createdAt,
			"updated_at":   updatedAt,
			"deleted_at":   deletedAt,
		}
		if locationID.Valid {
			doc["location_id"] = locationID.Int64
		} else {
			doc["location_id"] = nil
		}
		documents = append(documents, doc)
	}

	if documents == nil {
		documents = []map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

func (h *DocumentHandler) Restore(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.DB.Exec(`
		UPDATE documents d
		JOIN parties p ON d.party_id = p.id
		SET d.deleted_at = NULL
		WHERE d.id = ? AND p.user_id = ? AND d.deleted_at IS NOT NULL
	`, docID, userID)
	if err != nil {
		h.Logger.Error("failed to restore document", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"document not found in trash"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "restored"})
}

func (h *DocumentHandler) PermanentDelete(w http.ResponseWriter, r *http.Request) {
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

	// Verify ownership and that document is in trash
	var count int
	err = h.DB.QueryRow(`
		SELECT COUNT(*) FROM documents d
		JOIN parties p ON d.party_id = p.id
		WHERE d.id = ? AND p.user_id = ? AND d.deleted_at IS NOT NULL
	`, docID, userID).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, `{"error":"document not found in trash"}`, http.StatusNotFound)
		return
	}

	// Collect file paths before deleting
	fileRows, err := h.DB.Query("SELECT file_path FROM document_files WHERE document_id = ?", docID)
	if err != nil {
		h.Logger.Error("failed to query document files", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	var filePaths []string
	for fileRows.Next() {
		var path string
		if err := fileRows.Scan(&path); err != nil {
			continue
		}
		filePaths = append(filePaths, path)
	}
	fileRows.Close()

	// Delete file records
	_, _ = h.DB.Exec("DELETE FROM document_files WHERE document_id = ?", docID)

	// Hard delete the document
	_, err = h.DB.Exec(`
		DELETE d FROM documents d
		JOIN parties p ON d.party_id = p.id
		WHERE d.id = ? AND p.user_id = ? AND d.deleted_at IS NOT NULL
	`, docID, userID)
	if err != nil {
		h.Logger.Error("failed to permanently delete document", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Remove files from disk
	for _, path := range filePaths {
		os.Remove(path)
	}

	w.WriteHeader(http.StatusNoContent)
}
