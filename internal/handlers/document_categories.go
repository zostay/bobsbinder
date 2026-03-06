package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type DocumentCategoryHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func (h *DocumentCategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
		SELECT id, slug, name, description, sort_order
		FROM document_categories
		ORDER BY sort_order
	`)
	if err != nil {
		h.Logger.Error("failed to list document categories", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []map[string]any
	for rows.Next() {
		var id, sortOrder int64
		var slug, name string
		var description sql.NullString
		if err := rows.Scan(&id, &slug, &name, &description, &sortOrder); err != nil {
			h.Logger.Error("failed to scan document category", zap.Error(err))
			continue
		}
		cat := map[string]any{
			"id":         id,
			"slug":       slug,
			"name":       name,
			"sort_order": sortOrder,
		}
		if description.Valid {
			cat["description"] = description.String
		} else {
			cat["description"] = nil
		}
		categories = append(categories, cat)
	}

	if categories == nil {
		categories = []map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
