package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type PartyHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
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
