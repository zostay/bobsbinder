package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/zostay/bobsbinder/internal/services"
)

type SurvivorLetterHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type letterBoilerplateRequest struct {
	Greeting  string `json:"greeting"`
	Intro     string `json:"intro"`
	Closing   string `json:"closing"`
	Signature string `json:"signature"`
}

type sectionUpdateRequest struct {
	Title     *string `json:"title"`
	Visible   *bool   `json:"visible"`
	SortOrder *int    `json:"sort_order"`
}

type sectionReorderRequest struct {
	SectionOrders []struct {
		ID        int64 `json:"id"`
		SortOrder int   `json:"sort_order"`
	} `json:"section_orders"`
}

type addItemRequest struct {
	Content  string `json:"content"`
	ItemType string `json:"item_type"`
}

type editItemRequest struct {
	Content string `json:"content"`
}

type itemReorderRequest struct {
	ItemOrders []struct {
		ID        int64 `json:"id"`
		SortOrder int   `json:"sort_order"`
	} `json:"item_orders"`
}

// GetLetter runs sync and returns the full letter with all sections and items.
func (h *SurvivorLetterHandler) GetLetter(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	letter, err := services.SyncLetter(h.DB, h.Logger, userID)
	if err != nil {
		h.Logger.Error("failed to sync letter", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(letter)
}

// UpdateBoilerplate updates the greeting, intro, closing, and signature.
func (h *SurvivorLetterHandler) UpdateBoilerplate(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req letterBoilerplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE survivor_letters SET greeting = ?, intro = ?, closing = ?, signature = ?
		WHERE user_id = ?
	`, req.Greeting, req.Intro, req.Closing, req.Signature, userID)
	if err != nil {
		h.Logger.Error("failed to update letter boilerplate", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"letter not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// UpdateSection updates a section's title, visibility, or sort_order.
func (h *SurvivorLetterHandler) UpdateSection(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	sectionID, err := strconv.ParseInt(chi.URLParam(r, "sectionId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid section id"}`, http.StatusBadRequest)
		return
	}

	var req sectionUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Verify ownership through letter -> user
	var count int
	err = h.DB.QueryRow(`
		SELECT COUNT(*) FROM survivor_letter_sections sls
		JOIN survivor_letters sl ON sls.letter_id = sl.id
		WHERE sls.id = ? AND sl.user_id = ?
	`, sectionID, userID).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, `{"error":"section not found"}`, http.StatusNotFound)
		return
	}

	if req.Title != nil {
		h.DB.Exec("UPDATE survivor_letter_sections SET title = ? WHERE id = ?", *req.Title, sectionID)
	}
	if req.Visible != nil {
		h.DB.Exec("UPDATE survivor_letter_sections SET visible = ? WHERE id = ?", *req.Visible, sectionID)
	}
	if req.SortOrder != nil {
		h.DB.Exec("UPDATE survivor_letter_sections SET sort_order = ? WHERE id = ?", *req.SortOrder, sectionID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// ReorderSections batch-reorders sections.
func (h *SurvivorLetterHandler) ReorderSections(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req sectionReorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	for _, so := range req.SectionOrders {
		h.DB.Exec(`
			UPDATE survivor_letter_sections sls
			JOIN survivor_letters sl ON sls.letter_id = sl.id
			SET sls.sort_order = ?
			WHERE sls.id = ? AND sl.user_id = ?
		`, so.SortOrder, so.ID, userID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// AddItem adds a manual item to a section.
func (h *SurvivorLetterHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	sectionID, err := strconv.ParseInt(chi.URLParam(r, "sectionId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid section id"}`, http.StatusBadRequest)
		return
	}

	// Verify ownership
	var count int
	err = h.DB.QueryRow(`
		SELECT COUNT(*) FROM survivor_letter_sections sls
		JOIN survivor_letters sl ON sls.letter_id = sl.id
		WHERE sls.id = ? AND sl.user_id = ?
	`, sectionID, userID).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, `{"error":"section not found"}`, http.StatusNotFound)
		return
	}

	var req addItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.ItemType == "" {
		req.ItemType = "numbered"
	}

	var maxSort int
	h.DB.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM survivor_letter_items WHERE section_id = ?", sectionID).Scan(&maxSort)

	result, err := h.DB.Exec(`
		INSERT INTO survivor_letter_items (section_id, content, item_type, provenance, sort_order)
		VALUES (?, ?, ?, 'manual', ?)
	`, sectionID, req.Content, req.ItemType, maxSort+1)
	if err != nil {
		h.Logger.Error("failed to add manual item", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id": id, "section_id": sectionID, "content": req.Content,
		"item_type": req.ItemType, "provenance": "manual",
	})
}

// EditItem edits an item's content. Auto items become auto_edited.
func (h *SurvivorLetterHandler) EditItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid item id"}`, http.StatusBadRequest)
		return
	}

	var req editItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Verify ownership and get current provenance
	var provenance string
	err = h.DB.QueryRow(`
		SELECT sli.provenance FROM survivor_letter_items sli
		JOIN survivor_letter_sections sls ON sli.section_id = sls.id
		JOIN survivor_letters sl ON sls.letter_id = sl.id
		WHERE sli.id = ? AND sl.user_id = ?
	`, itemID, userID).Scan(&provenance)
	if err != nil {
		http.Error(w, `{"error":"item not found"}`, http.StatusNotFound)
		return
	}

	newProvenance := provenance
	if provenance == "auto" {
		newProvenance = "auto_edited"
	}

	_, err = h.DB.Exec("UPDATE survivor_letter_items SET content = ?, provenance = ? WHERE id = ?",
		req.Content, newProvenance, itemID)
	if err != nil {
		h.Logger.Error("failed to edit item", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"status": "updated", "provenance": newProvenance})
}

// ReorderItems batch-reorders items.
func (h *SurvivorLetterHandler) ReorderItems(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req itemReorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	for _, io := range req.ItemOrders {
		h.DB.Exec(`
			UPDATE survivor_letter_items sli
			JOIN survivor_letter_sections sls ON sli.section_id = sls.id
			JOIN survivor_letters sl ON sls.letter_id = sl.id
			SET sli.sort_order = ?
			WHERE sli.id = ? AND sl.user_id = ?
		`, io.SortOrder, io.ID, userID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// DeleteItem hard-deletes manual items, or suppresses auto items.
func (h *SurvivorLetterHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid item id"}`, http.StatusBadRequest)
		return
	}

	var provenance string
	err = h.DB.QueryRow(`
		SELECT sli.provenance FROM survivor_letter_items sli
		JOIN survivor_letter_sections sls ON sli.section_id = sls.id
		JOIN survivor_letters sl ON sls.letter_id = sl.id
		WHERE sli.id = ? AND sl.user_id = ?
	`, itemID, userID).Scan(&provenance)
	if err != nil {
		http.Error(w, `{"error":"item not found"}`, http.StatusNotFound)
		return
	}

	if provenance == "manual" {
		_, err = h.DB.Exec("DELETE FROM survivor_letter_items WHERE id = ?", itemID)
	} else {
		_, err = h.DB.Exec("UPDATE survivor_letter_items SET suppressed = TRUE WHERE id = ?", itemID)
	}
	if err != nil {
		h.Logger.Error("failed to delete/suppress item", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnsuppressItem restores a suppressed auto item.
func (h *SurvivorLetterHandler) UnsuppressItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemId"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid item id"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE survivor_letter_items sli
		JOIN survivor_letter_sections sls ON sli.section_id = sls.id
		JOIN survivor_letters sl ON sls.letter_id = sl.id
		SET sli.suppressed = FALSE
		WHERE sli.id = ? AND sl.user_id = ?
	`, itemID, userID)
	if err != nil {
		h.Logger.Error("failed to unsuppress item", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"item not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}
