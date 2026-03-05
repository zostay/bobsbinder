package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type DocumentHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type documentRequest struct {
	PartyID    int64  `json:"party_id"`
	CategoryID int64  `json:"category_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Status     string `json:"status"`
}

func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(`
		SELECT d.id, d.party_id, d.category_id, d.title, d.content, d.status, d.created_at, d.updated_at
		FROM documents d
		JOIN parties p ON d.party_id = p.id
		WHERE p.user_id = ?
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
		var title, content, status string
		var createdAt, updatedAt string
		if err := rows.Scan(&id, &partyID, &categoryID, &title, &content, &status, &createdAt, &updatedAt); err != nil {
			h.Logger.Error("failed to scan document", zap.Error(err))
			continue
		}
		documents = append(documents, map[string]any{
			"id":          id,
			"party_id":    partyID,
			"category_id": categoryID,
			"title":       title,
			"content":     content,
			"status":      status,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		})
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
	var title, content, status string
	var createdAt, updatedAt string
	err = h.DB.QueryRow(`
		SELECT d.id, d.party_id, d.category_id, d.title, d.content, d.status, d.created_at, d.updated_at
		FROM documents d
		JOIN parties p ON d.party_id = p.id
		WHERE d.id = ? AND p.user_id = ?
	`, docID, userID).Scan(&id, &partyID, &categoryID, &title, &content, &status, &createdAt, &updatedAt)
	if err != nil {
		http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":          id,
		"party_id":    partyID,
		"category_id": categoryID,
		"title":       title,
		"content":     content,
		"status":      status,
		"created_at":  createdAt,
		"updated_at":  updatedAt,
	})
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

	result, err := h.DB.Exec(
		"INSERT INTO documents (party_id, category_id, title, content, status) VALUES (?, ?, ?, ?, ?)",
		req.PartyID, req.CategoryID, req.Title, req.Content, req.Status,
	)
	if err != nil {
		h.Logger.Error("failed to create document", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	docID, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":          docID,
		"party_id":    req.PartyID,
		"category_id": req.CategoryID,
		"title":       req.Title,
		"content":     req.Content,
		"status":      req.Status,
	})
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
		SET d.title = ?, d.content = ?, d.status = ?
		WHERE d.id = ? AND p.user_id = ?
	`, req.Title, req.Content, req.Status, docID, userID)
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
		DELETE d FROM documents d
		JOIN parties p ON d.party_id = p.id
		WHERE d.id = ? AND p.user_id = ?
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
