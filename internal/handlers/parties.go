package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type PartyHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type partyRequest struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Notes        string `json:"notes"`
}

var validRelationships = map[string]bool{
	"self": true, "spouse": true, "dependent": true, "other": true,
}

func (h *PartyHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(`
		SELECT id, user_id, name, relationship, notes, created_at, updated_at
		FROM parties WHERE user_id = ?
		ORDER BY (relationship = 'self') DESC, name
	`, userID)
	if err != nil {
		h.Logger.Error("failed to list parties", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var parties []map[string]any
	for rows.Next() {
		var id, uid int64
		var name, relationship, createdAt, updatedAt string
		var notes sql.NullString
		if err := rows.Scan(&id, &uid, &name, &relationship, &notes, &createdAt, &updatedAt); err != nil {
			h.Logger.Error("failed to scan party", zap.Error(err))
			continue
		}
		party := map[string]any{
			"id":           id,
			"user_id":      uid,
			"name":         name,
			"relationship": relationship,
			"created_at":   createdAt,
			"updated_at":   updatedAt,
		}
		if notes.Valid {
			party["notes"] = notes.String
		} else {
			party["notes"] = nil
		}
		parties = append(parties, party)
	}

	if parties == nil {
		parties = []map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parties)
}

func (h *PartyHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req partyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if !validRelationships[req.Relationship] {
		http.Error(w, `{"error":"invalid relationship, must be one of: self, spouse, dependent, other"}`, http.StatusBadRequest)
		return
	}

	if req.Relationship == "self" {
		var count int
		h.DB.QueryRow("SELECT COUNT(*) FROM parties WHERE user_id = ? AND relationship = 'self'", userID).Scan(&count)
		if count > 0 {
			http.Error(w, `{"error":"a self party already exists"}`, http.StatusConflict)
			return
		}
	}

	result, err := h.DB.Exec(
		"INSERT INTO parties (user_id, name, relationship, notes) VALUES (?, ?, ?, ?)",
		userID, req.Name, req.Relationship, req.Notes,
	)
	if err != nil {
		h.Logger.Error("failed to create party", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id": id, "user_id": userID, "name": req.Name,
		"relationship": req.Relationship, "notes": req.Notes,
	})
}

func (h *PartyHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	partyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid party id"}`, http.StatusBadRequest)
		return
	}

	var req partyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if !validRelationships[req.Relationship] {
		http.Error(w, `{"error":"invalid relationship"}`, http.StatusBadRequest)
		return
	}

	if req.Relationship == "self" {
		var existingID int64
		err := h.DB.QueryRow("SELECT id FROM parties WHERE user_id = ? AND relationship = 'self'", userID).Scan(&existingID)
		if err == nil && existingID != partyID {
			http.Error(w, `{"error":"a self party already exists"}`, http.StatusConflict)
			return
		}
	}

	result, err := h.DB.Exec(`
		UPDATE parties SET name = ?, relationship = ?, notes = ?
		WHERE id = ? AND user_id = ?
	`, req.Name, req.Relationship, req.Notes, partyID, userID)
	if err != nil {
		h.Logger.Error("failed to update party", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"party not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *PartyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	partyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid party id"}`, http.StatusBadRequest)
		return
	}

	// Prevent deleting self party
	var relationship string
	err = h.DB.QueryRow("SELECT relationship FROM parties WHERE id = ? AND user_id = ?", partyID, userID).Scan(&relationship)
	if err != nil {
		http.Error(w, `{"error":"party not found"}`, http.StatusNotFound)
		return
	}
	if relationship == "self" {
		http.Error(w, `{"error":"cannot delete self party"}`, http.StatusForbidden)
		return
	}

	result, err := h.DB.Exec("DELETE FROM parties WHERE id = ? AND user_id = ?", partyID, userID)
	if err != nil {
		h.Logger.Error("failed to delete party", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"party not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
