package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type DigitalAccessHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type digitalAccessRequest struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Instructions string `json:"instructions"`
	LocationID   *int64 `json:"location_id"`
	SecureNotes  string `json:"secure_notes"`
}

func (h *DigitalAccessHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(`
		SELECT id, user_id, type, name, username, instructions, secure_notes, location_id, created_at, updated_at
		FROM digital_access WHERE user_id = ? ORDER BY name
	`, userID)
	if err != nil {
		h.Logger.Error("failed to list digital access", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []map[string]any
	for rows.Next() {
		var id, uid int64
		var daType, name, username, instructions, secureNotes, createdAt, updatedAt string
		var locationID sql.NullInt64
		if err := rows.Scan(&id, &uid, &daType, &name, &username, &instructions, &secureNotes, &locationID, &createdAt, &updatedAt); err != nil {
			h.Logger.Error("failed to scan digital access", zap.Error(err))
			continue
		}
		item := map[string]any{
			"id": id, "user_id": uid, "type": daType, "name": name,
			"username": username, "instructions": instructions, "secure_notes": secureNotes,
			"created_at": createdAt, "updated_at": updatedAt,
		}
		if locationID.Valid {
			item["location_id"] = locationID.Int64
		} else {
			item["location_id"] = nil
		}
		items = append(items, item)
	}

	if items == nil {
		items = []map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *DigitalAccessHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	daID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid digital access id"}`, http.StatusBadRequest)
		return
	}

	var id, uid int64
	var daType, name, username, instructions, secureNotes, createdAt, updatedAt string
	var locationID sql.NullInt64
	err = h.DB.QueryRow(`
		SELECT id, user_id, type, name, username, instructions, secure_notes, location_id, created_at, updated_at
		FROM digital_access WHERE id = ? AND user_id = ?
	`, daID, userID).Scan(&id, &uid, &daType, &name, &username, &instructions, &secureNotes, &locationID, &createdAt, &updatedAt)
	if err != nil {
		http.Error(w, `{"error":"digital access not found"}`, http.StatusNotFound)
		return
	}

	item := map[string]any{
		"id": id, "user_id": uid, "type": daType, "name": name,
		"username": username, "instructions": instructions, "secure_notes": secureNotes,
		"created_at": createdAt, "updated_at": updatedAt,
	}
	if locationID.Valid {
		item["location_id"] = locationID.Int64
	} else {
		item["location_id"] = nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *DigitalAccessHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req digitalAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(
		"INSERT INTO digital_access (user_id, type, name, username, instructions, secure_notes, location_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID, req.Type, req.Name, req.Username, req.Instructions, req.SecureNotes, req.LocationID,
	)
	if err != nil {
		h.Logger.Error("failed to create digital access", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id": id, "user_id": userID, "type": req.Type, "name": req.Name,
		"username": req.Username, "instructions": req.Instructions, "secure_notes": req.SecureNotes,
		"location_id": req.LocationID,
	})
}

func (h *DigitalAccessHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	daID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid digital access id"}`, http.StatusBadRequest)
		return
	}

	var req digitalAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE digital_access SET type = ?, name = ?, username = ?, instructions = ?, secure_notes = ?, location_id = ?
		WHERE id = ? AND user_id = ?
	`, req.Type, req.Name, req.Username, req.Instructions, req.SecureNotes, req.LocationID, daID, userID)
	if err != nil {
		h.Logger.Error("failed to update digital access", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"digital access not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *DigitalAccessHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	daID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid digital access id"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec("DELETE FROM digital_access WHERE id = ? AND user_id = ?", daID, userID)
	if err != nil {
		h.Logger.Error("failed to delete digital access", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"digital access not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
