package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ChecklistHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type checklistCategory struct {
	CategoryID   int64  `json:"category_id"`
	CategorySlug string `json:"category_slug"`
	CategoryName string `json:"category_name"`
	Status       string `json:"status"`
	HasDocument  bool   `json:"has_document"`
}

type partyChecklist struct {
	PartyID           int64               `json:"party_id"`
	PartyName         string              `json:"party_name"`
	PartyRelationship string              `json:"party_relationship"`
	Items             []checklistCategory `json:"items"`
}

type checklistStatusRequest struct {
	Status string `json:"status"`
}

var validStatuses = map[string]bool{
	"pending": true, "complete": true, "not_applicable": true,
}

func (h *ChecklistHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	checklists, err := h.buildChecklists(userID, 0)
	if err != nil {
		h.Logger.Error("failed to build checklists", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checklists)
}

func (h *ChecklistHandler) ListForParty(w http.ResponseWriter, r *http.Request) {
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

	// Verify ownership
	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM parties WHERE id = ? AND user_id = ?", partyID, userID).Scan(&count)
	if count == 0 {
		http.Error(w, `{"error":"party not found"}`, http.StatusNotFound)
		return
	}

	checklists, err := h.buildChecklists(userID, partyID)
	if err != nil {
		h.Logger.Error("failed to build party checklist", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checklists)
}

func (h *ChecklistHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
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

	categoryID, err := strconv.ParseInt(chi.URLParam(r, "categoryId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid category id"}`, http.StatusBadRequest)
		return
	}

	// Verify party ownership
	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM parties WHERE id = ? AND user_id = ?", partyID, userID).Scan(&count)
	if count == 0 {
		http.Error(w, `{"error":"party not found"}`, http.StatusNotFound)
		return
	}

	var req checklistStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if !validStatuses[req.Status] {
		http.Error(w, `{"error":"invalid status, must be one of: pending, complete, not_applicable"}`, http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec(`
		INSERT INTO checklist_items (party_id, category_id, status)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE status = VALUES(status)
	`, partyID, categoryID, req.Status)
	if err != nil {
		h.Logger.Error("failed to update checklist status", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": req.Status})
}

func (h *ChecklistHandler) buildChecklists(userID int64, filterPartyID int64) ([]partyChecklist, error) {
	// Get parties
	partyQuery := "SELECT id, name, relationship FROM parties WHERE user_id = ? ORDER BY (relationship = 'self') DESC, name"
	partyArgs := []any{userID}
	if filterPartyID > 0 {
		partyQuery = "SELECT id, name, relationship FROM parties WHERE id = ? AND user_id = ?"
		partyArgs = []any{filterPartyID, userID}
	}

	partyRows, err := h.DB.Query(partyQuery, partyArgs...)
	if err != nil {
		return nil, err
	}
	defer partyRows.Close()

	type partyInfo struct {
		ID           int64
		Name         string
		Relationship string
	}
	var parties []partyInfo
	for partyRows.Next() {
		var p partyInfo
		if err := partyRows.Scan(&p.ID, &p.Name, &p.Relationship); err != nil {
			continue
		}
		parties = append(parties, p)
	}

	var result []partyChecklist
	for _, p := range parties {
		rows, err := h.DB.Query(`
			SELECT dc.id, dc.slug, dc.name,
				COALESCE(ci.status, 'pending') as status,
				EXISTS(SELECT 1 FROM documents d WHERE d.party_id = ? AND d.category_id = dc.id AND d.deleted_at IS NULL) as has_document
			FROM document_categories dc
			LEFT JOIN checklist_items ci ON ci.category_id = dc.id AND ci.party_id = ?
			ORDER BY dc.sort_order
		`, p.ID, p.ID)
		if err != nil {
			return nil, err
		}

		var items []checklistCategory
		for rows.Next() {
			var item checklistCategory
			if err := rows.Scan(&item.CategoryID, &item.CategorySlug, &item.CategoryName, &item.Status, &item.HasDocument); err != nil {
				rows.Close()
				return nil, err
			}
			items = append(items, item)
		}
		rows.Close()

		if items == nil {
			items = []checklistCategory{}
		}

		result = append(result, partyChecklist{
			PartyID:           p.ID,
			PartyName:         p.Name,
			PartyRelationship: p.Relationship,
			Items:             items,
		})
	}

	if result == nil {
		result = []partyChecklist{}
	}

	return result, nil
}
