package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ServiceAccountHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type serviceAccountRequest struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	Provider      string `json:"provider"`
	AccountNumber string `json:"account_number"`
	ContactName   string `json:"contact_name"`
	ContactPhone  string `json:"contact_phone"`
	ContactEmail  string `json:"contact_email"`
	Notes         string `json:"notes"`
	SecureNotes   string `json:"secure_notes"`
}

func (h *ServiceAccountHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(`
		SELECT id, user_id, type, name, provider, account_number, contact_name, contact_phone, contact_email, notes, secure_notes, created_at, updated_at
		FROM service_accounts WHERE user_id = ? ORDER BY name
	`, userID)
	if err != nil {
		h.Logger.Error("failed to list service accounts", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var accounts []map[string]any
	for rows.Next() {
		var id, uid int64
		var saType, name, provider, accountNumber, contactName, contactPhone, contactEmail, notes, secureNotes, createdAt, updatedAt string
		if err := rows.Scan(&id, &uid, &saType, &name, &provider, &accountNumber, &contactName, &contactPhone, &contactEmail, &notes, &secureNotes, &createdAt, &updatedAt); err != nil {
			h.Logger.Error("failed to scan service account", zap.Error(err))
			continue
		}
		accounts = append(accounts, map[string]any{
			"id": id, "user_id": uid, "type": saType, "name": name,
			"provider": provider, "account_number": accountNumber,
			"contact_name": contactName, "contact_phone": contactPhone,
			"contact_email": contactEmail, "notes": notes, "secure_notes": secureNotes,
			"created_at": createdAt, "updated_at": updatedAt,
		})
	}

	if accounts == nil {
		accounts = []map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func (h *ServiceAccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	saID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid service account id"}`, http.StatusBadRequest)
		return
	}

	var id, uid int64
	var saType, name, provider, accountNumber, contactName, contactPhone, contactEmail, notes, secureNotes, createdAt, updatedAt string
	err = h.DB.QueryRow(`
		SELECT id, user_id, type, name, provider, account_number, contact_name, contact_phone, contact_email, notes, secure_notes, created_at, updated_at
		FROM service_accounts WHERE id = ? AND user_id = ?
	`, saID, userID).Scan(&id, &uid, &saType, &name, &provider, &accountNumber, &contactName, &contactPhone, &contactEmail, &notes, &secureNotes, &createdAt, &updatedAt)
	if err != nil {
		http.Error(w, `{"error":"service account not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id": id, "user_id": uid, "type": saType, "name": name,
		"provider": provider, "account_number": accountNumber,
		"contact_name": contactName, "contact_phone": contactPhone,
		"contact_email": contactEmail, "notes": notes, "secure_notes": secureNotes,
		"created_at": createdAt, "updated_at": updatedAt,
	})
}

func (h *ServiceAccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req serviceAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(
		"INSERT INTO service_accounts (user_id, type, name, provider, account_number, contact_name, contact_phone, contact_email, notes, secure_notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		userID, req.Type, req.Name, req.Provider, req.AccountNumber, req.ContactName, req.ContactPhone, req.ContactEmail, req.Notes, req.SecureNotes,
	)
	if err != nil {
		h.Logger.Error("failed to create service account", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id": id, "user_id": userID, "type": req.Type, "name": req.Name,
		"provider": req.Provider, "account_number": req.AccountNumber,
		"contact_name": req.ContactName, "contact_phone": req.ContactPhone,
		"contact_email": req.ContactEmail, "notes": req.Notes, "secure_notes": req.SecureNotes,
	})
}

func (h *ServiceAccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	saID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid service account id"}`, http.StatusBadRequest)
		return
	}

	var req serviceAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE service_accounts SET type = ?, name = ?, provider = ?, account_number = ?,
		       contact_name = ?, contact_phone = ?, contact_email = ?, notes = ?, secure_notes = ?
		WHERE id = ? AND user_id = ?
	`, req.Type, req.Name, req.Provider, req.AccountNumber, req.ContactName, req.ContactPhone, req.ContactEmail, req.Notes, req.SecureNotes, saID, userID)
	if err != nil {
		h.Logger.Error("failed to update service account", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"service account not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *ServiceAccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	saID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid service account id"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec("DELETE FROM service_accounts WHERE id = ? AND user_id = ?", saID, userID)
	if err != nil {
		h.Logger.Error("failed to delete service account", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"service account not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
