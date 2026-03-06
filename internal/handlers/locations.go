package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type LocationHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type locationRequest struct {
	Name               string `json:"name"`
	Type               string `json:"type"`
	Description        string `json:"description"`
	Address            string `json:"address"`
	AccessInstructions string `json:"access_instructions"`
	SecureNotes        string `json:"secure_notes"`
}

func (h *LocationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(`
		SELECT id, user_id, name, type, description, address, access_instructions, secure_notes, created_at, updated_at
		FROM locations WHERE user_id = ? ORDER BY name
	`, userID)
	if err != nil {
		h.Logger.Error("failed to list locations", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var locations []map[string]any
	for rows.Next() {
		var id, uid int64
		var name, locType, description, address, accessInstructions, secureNotes, createdAt, updatedAt string
		if err := rows.Scan(&id, &uid, &name, &locType, &description, &address, &accessInstructions, &secureNotes, &createdAt, &updatedAt); err != nil {
			h.Logger.Error("failed to scan location", zap.Error(err))
			continue
		}
		locations = append(locations, map[string]any{
			"id": id, "user_id": uid, "name": name, "type": locType,
			"description": description, "address": address, "access_instructions": accessInstructions,
			"secure_notes": secureNotes, "created_at": createdAt, "updated_at": updatedAt,
		})
	}

	if locations == nil {
		locations = []map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func (h *LocationHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	locID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid location id"}`, http.StatusBadRequest)
		return
	}

	var id, uid int64
	var name, locType, description, address, accessInstructions, secureNotes, createdAt, updatedAt string
	err = h.DB.QueryRow(`
		SELECT id, user_id, name, type, description, address, access_instructions, secure_notes, created_at, updated_at
		FROM locations WHERE id = ? AND user_id = ?
	`, locID, userID).Scan(&id, &uid, &name, &locType, &description, &address, &accessInstructions, &secureNotes, &createdAt, &updatedAt)
	if err != nil {
		http.Error(w, `{"error":"location not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id": id, "user_id": uid, "name": name, "type": locType,
		"description": description, "address": address, "access_instructions": accessInstructions,
		"secure_notes": secureNotes, "created_at": createdAt, "updated_at": updatedAt,
	})
}

func (h *LocationHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req locationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(
		"INSERT INTO locations (user_id, name, type, description, address, access_instructions, secure_notes) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID, req.Name, req.Type, req.Description, req.Address, req.AccessInstructions, req.SecureNotes,
	)
	if err != nil {
		h.Logger.Error("failed to create location", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id": id, "user_id": userID, "name": req.Name, "type": req.Type,
		"description": req.Description, "address": req.Address, "access_instructions": req.AccessInstructions,
		"secure_notes": req.SecureNotes,
	})
}

func (h *LocationHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	locID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid location id"}`, http.StatusBadRequest)
		return
	}

	var req locationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE locations SET name = ?, type = ?, description = ?, address = ?, access_instructions = ?, secure_notes = ?
		WHERE id = ? AND user_id = ?
	`, req.Name, req.Type, req.Description, req.Address, req.AccessInstructions, req.SecureNotes, locID, userID)
	if err != nil {
		h.Logger.Error("failed to update location", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"location not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *LocationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	locID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid location id"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec("DELETE FROM locations WHERE id = ? AND user_id = ?", locID, userID)
	if err != nil {
		h.Logger.Error("failed to delete location", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"location not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
