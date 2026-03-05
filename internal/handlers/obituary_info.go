package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ObituaryInfoHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type obituaryInfoRequest struct {
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	Relationship string  `json:"relationship"`
	Details      string  `json:"details"`
	EventDate    *string `json:"event_date"`
}

func (h *ObituaryInfoHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	partyID, err := strconv.ParseInt(chi.URLParam(r, "partyId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid party id"}`, http.StatusBadRequest)
		return
	}

	// Verify party belongs to user
	var count int
	if err := h.DB.QueryRow("SELECT COUNT(*) FROM parties WHERE id = ? AND user_id = ?", partyID, userID).Scan(&count); err != nil || count == 0 {
		http.Error(w, `{"error":"party not found"}`, http.StatusNotFound)
		return
	}

	rows, err := h.DB.Query(`
		SELECT id, party_id, type, name, relationship, details, event_date, created_at, updated_at
		FROM party_obituary_info WHERE party_id = ? ORDER BY name
	`, partyID)
	if err != nil {
		h.Logger.Error("failed to list obituary info", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []map[string]any
	for rows.Next() {
		var id, pid int64
		var oType, name, relationship, details, createdAt, updatedAt string
		var eventDate sql.NullString
		if err := rows.Scan(&id, &pid, &oType, &name, &relationship, &details, &eventDate, &createdAt, &updatedAt); err != nil {
			h.Logger.Error("failed to scan obituary info", zap.Error(err))
			continue
		}
		item := map[string]any{
			"id": id, "party_id": pid, "type": oType, "name": name,
			"relationship": relationship, "details": details,
			"created_at": createdAt, "updated_at": updatedAt,
		}
		if eventDate.Valid {
			item["event_date"] = eventDate.String
		} else {
			item["event_date"] = nil
		}
		items = append(items, item)
	}

	if items == nil {
		items = []map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *ObituaryInfoHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	partyID, err := strconv.ParseInt(chi.URLParam(r, "partyId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid party id"}`, http.StatusBadRequest)
		return
	}

	var count int
	if err := h.DB.QueryRow("SELECT COUNT(*) FROM parties WHERE id = ? AND user_id = ?", partyID, userID).Scan(&count); err != nil || count == 0 {
		http.Error(w, `{"error":"party not found"}`, http.StatusNotFound)
		return
	}

	var req obituaryInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(
		"INSERT INTO party_obituary_info (party_id, type, name, relationship, details, event_date) VALUES (?, ?, ?, ?, ?, ?)",
		partyID, req.Type, req.Name, req.Relationship, req.Details, req.EventDate,
	)
	if err != nil {
		h.Logger.Error("failed to create obituary info", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id": id, "party_id": partyID, "type": req.Type, "name": req.Name,
		"relationship": req.Relationship, "details": req.Details, "event_date": req.EventDate,
	})
}

func (h *ObituaryInfoHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	partyID, err := strconv.ParseInt(chi.URLParam(r, "partyId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid party id"}`, http.StatusBadRequest)
		return
	}

	infoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid obituary info id"}`, http.StatusBadRequest)
		return
	}

	var count int
	if err := h.DB.QueryRow("SELECT COUNT(*) FROM parties WHERE id = ? AND user_id = ?", partyID, userID).Scan(&count); err != nil || count == 0 {
		http.Error(w, `{"error":"party not found"}`, http.StatusNotFound)
		return
	}

	var req obituaryInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE party_obituary_info SET type = ?, name = ?, relationship = ?, details = ?, event_date = ?
		WHERE id = ? AND party_id = ?
	`, req.Type, req.Name, req.Relationship, req.Details, req.EventDate, infoID, partyID)
	if err != nil {
		h.Logger.Error("failed to update obituary info", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"obituary info not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *ObituaryInfoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	partyID, err := strconv.ParseInt(chi.URLParam(r, "partyId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid party id"}`, http.StatusBadRequest)
		return
	}

	infoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid obituary info id"}`, http.StatusBadRequest)
		return
	}

	var count int
	if err := h.DB.QueryRow("SELECT COUNT(*) FROM parties WHERE id = ? AND user_id = ?", partyID, userID).Scan(&count); err != nil || count == 0 {
		http.Error(w, `{"error":"party not found"}`, http.StatusNotFound)
		return
	}

	result, err := h.DB.Exec("DELETE FROM party_obituary_info WHERE id = ? AND party_id = ?", infoID, partyID)
	if err != nil {
		h.Logger.Error("failed to delete obituary info", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"obituary info not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
