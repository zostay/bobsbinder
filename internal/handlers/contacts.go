package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ContactHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type contactRequest struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Role         string `json:"role"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Notes        string `json:"notes"`
	IsPrimary    bool   `json:"is_primary"`
}

func (h *ContactHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(`
		SELECT id, user_id, name, relationship, role, phone, email, address, notes, is_primary, created_at, updated_at
		FROM contacts WHERE user_id = ? ORDER BY is_primary DESC, name
	`, userID)
	if err != nil {
		h.Logger.Error("failed to list contacts", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var contacts []map[string]any
	for rows.Next() {
		var id, uid int64
		var name, relationship, role, phone, email, address, notes, createdAt, updatedAt string
		var isPrimary bool
		if err := rows.Scan(&id, &uid, &name, &relationship, &role, &phone, &email, &address, &notes, &isPrimary, &createdAt, &updatedAt); err != nil {
			h.Logger.Error("failed to scan contact", zap.Error(err))
			continue
		}
		contacts = append(contacts, map[string]any{
			"id": id, "user_id": uid, "name": name, "relationship": relationship,
			"role": role, "phone": phone, "email": email, "address": address,
			"notes": notes, "is_primary": isPrimary, "created_at": createdAt, "updated_at": updatedAt,
		})
	}

	if contacts == nil {
		contacts = []map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

func (h *ContactHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	contactID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid contact id"}`, http.StatusBadRequest)
		return
	}

	var id, uid int64
	var name, relationship, role, phone, email, address, notes, createdAt, updatedAt string
	var isPrimary bool
	err = h.DB.QueryRow(`
		SELECT id, user_id, name, relationship, role, phone, email, address, notes, is_primary, created_at, updated_at
		FROM contacts WHERE id = ? AND user_id = ?
	`, contactID, userID).Scan(&id, &uid, &name, &relationship, &role, &phone, &email, &address, &notes, &isPrimary, &createdAt, &updatedAt)
	if err != nil {
		http.Error(w, `{"error":"contact not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id": id, "user_id": uid, "name": name, "relationship": relationship,
		"role": role, "phone": phone, "email": email, "address": address,
		"notes": notes, "is_primary": isPrimary, "created_at": createdAt, "updated_at": updatedAt,
	})
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req contactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(
		"INSERT INTO contacts (user_id, name, relationship, role, phone, email, address, notes, is_primary) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		userID, req.Name, req.Relationship, req.Role, req.Phone, req.Email, req.Address, req.Notes, req.IsPrimary,
	)
	if err != nil {
		h.Logger.Error("failed to create contact", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id": id, "user_id": userID, "name": req.Name, "relationship": req.Relationship,
		"role": req.Role, "phone": req.Phone, "email": req.Email, "address": req.Address,
		"notes": req.Notes, "is_primary": req.IsPrimary,
	})
}

func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	contactID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid contact id"}`, http.StatusBadRequest)
		return
	}

	var req contactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE contacts SET name = ?, relationship = ?, role = ?, phone = ?, email = ?, address = ?, notes = ?, is_primary = ?
		WHERE id = ? AND user_id = ?
	`, req.Name, req.Relationship, req.Role, req.Phone, req.Email, req.Address, req.Notes, req.IsPrimary, contactID, userID)
	if err != nil {
		h.Logger.Error("failed to update contact", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"contact not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	contactID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid contact id"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec("DELETE FROM contacts WHERE id = ? AND user_id = ?", contactID, userID)
	if err != nil {
		h.Logger.Error("failed to delete contact", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"contact not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
