package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type ConfidentialHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type confidentialField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type confidentialItem struct {
	Label  string              `json:"label"`
	Fields []confidentialField `json:"fields"`
}

type confidentialSection struct {
	Title string             `json:"title"`
	Items []confidentialItem `json:"items"`
}

func (h *ConfidentialHandler) GetConfidential(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var sections []confidentialSection

	// Contacts
	if items, err := h.queryContacts(userID); err == nil && len(items) > 0 {
		sections = append(sections, confidentialSection{Title: "Contacts", Items: items})
	}

	// Documents
	if items, err := h.queryDocuments(userID); err == nil && len(items) > 0 {
		sections = append(sections, confidentialSection{Title: "Documents", Items: items})
	}

	// Insurance Policies
	if items, err := h.queryInsurancePolicies(userID); err == nil && len(items) > 0 {
		sections = append(sections, confidentialSection{Title: "Insurance Policies", Items: items})
	}

	// Important Locations
	if items, err := h.queryLocations(userID); err == nil && len(items) > 0 {
		sections = append(sections, confidentialSection{Title: "Important Locations", Items: items})
	}

	// Digital Information (digital_access + service_accounts combined)
	var digitalItems []confidentialItem
	if items, err := h.queryDigitalAccess(userID); err == nil {
		digitalItems = append(digitalItems, items...)
	}
	if items, err := h.queryServiceAccounts(userID); err == nil {
		digitalItems = append(digitalItems, items...)
	}
	if len(digitalItems) > 0 {
		sections = append(sections, confidentialSection{Title: "Digital Information", Items: digitalItems})
	}

	if sections == nil {
		sections = []confidentialSection{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sections)
}

func (h *ConfidentialHandler) queryContacts(userID int64) ([]confidentialItem, error) {
	rows, err := h.DB.Query(`
		SELECT name, secure_notes FROM contacts
		WHERE user_id = ? AND secure_notes != ''
		ORDER BY name
	`, userID)
	if err != nil {
		h.Logger.Error("failed to query confidential contacts", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var items []confidentialItem
	for rows.Next() {
		var name, secureNotes string
		if err := rows.Scan(&name, &secureNotes); err != nil {
			continue
		}
		var fields []confidentialField
		if secureNotes != "" {
			fields = append(fields, confidentialField{Name: "Confidential Notes", Value: secureNotes})
		}
		items = append(items, confidentialItem{Label: name, Fields: fields})
	}
	return items, nil
}

func (h *ConfidentialHandler) queryDocuments(userID int64) ([]confidentialItem, error) {
	rows, err := h.DB.Query(`
		SELECT d.title, d.secure_notes FROM documents d
		JOIN parties p ON d.party_id = p.id
		WHERE p.user_id = ? AND d.secure_notes != '' AND d.deleted_at IS NULL
		ORDER BY d.title
	`, userID)
	if err != nil {
		h.Logger.Error("failed to query confidential documents", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var items []confidentialItem
	for rows.Next() {
		var title, secureNotes string
		if err := rows.Scan(&title, &secureNotes); err != nil {
			continue
		}
		var fields []confidentialField
		if secureNotes != "" {
			fields = append(fields, confidentialField{Name: "Confidential Notes", Value: secureNotes})
		}
		items = append(items, confidentialItem{Label: title, Fields: fields})
	}
	return items, nil
}

func (h *ConfidentialHandler) queryInsurancePolicies(userID int64) ([]confidentialItem, error) {
	rows, err := h.DB.Query(`
		SELECT provider, type, policy_number, secure_notes FROM insurance_policies
		WHERE user_id = ? AND (policy_number != '' OR secure_notes != '')
		ORDER BY provider
	`, userID)
	if err != nil {
		h.Logger.Error("failed to query confidential insurance policies", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var items []confidentialItem
	for rows.Next() {
		var provider, pType, policyNumber, secureNotes string
		if err := rows.Scan(&provider, &pType, &policyNumber, &secureNotes); err != nil {
			continue
		}
		label := provider
		if pType != "" {
			label += fmt.Sprintf(" (%s)", pType)
		}
		var fields []confidentialField
		if policyNumber != "" {
			fields = append(fields, confidentialField{Name: "Policy Number", Value: policyNumber})
		}
		if secureNotes != "" {
			fields = append(fields, confidentialField{Name: "Confidential Notes", Value: secureNotes})
		}
		items = append(items, confidentialItem{Label: label, Fields: fields})
	}
	return items, nil
}

func (h *ConfidentialHandler) queryLocations(userID int64) ([]confidentialItem, error) {
	rows, err := h.DB.Query(`
		SELECT name, access_instructions, secure_notes FROM locations
		WHERE user_id = ? AND (access_instructions != '' OR secure_notes != '')
		ORDER BY name
	`, userID)
	if err != nil {
		h.Logger.Error("failed to query confidential locations", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var items []confidentialItem
	for rows.Next() {
		var name, accessInstructions, secureNotes string
		if err := rows.Scan(&name, &accessInstructions, &secureNotes); err != nil {
			continue
		}
		var fields []confidentialField
		if accessInstructions != "" {
			fields = append(fields, confidentialField{Name: "Access Instructions", Value: accessInstructions})
		}
		if secureNotes != "" {
			fields = append(fields, confidentialField{Name: "Confidential Notes", Value: secureNotes})
		}
		items = append(items, confidentialItem{Label: name, Fields: fields})
	}
	return items, nil
}

func (h *ConfidentialHandler) queryDigitalAccess(userID int64) ([]confidentialItem, error) {
	rows, err := h.DB.Query(`
		SELECT name, username, instructions, secure_notes FROM digital_access
		WHERE user_id = ? AND (username != '' OR instructions != '' OR secure_notes != '')
		ORDER BY name
	`, userID)
	if err != nil {
		h.Logger.Error("failed to query confidential digital access", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var items []confidentialItem
	for rows.Next() {
		var name, username, instructions, secureNotes string
		if err := rows.Scan(&name, &username, &instructions, &secureNotes); err != nil {
			continue
		}
		var fields []confidentialField
		if username != "" {
			fields = append(fields, confidentialField{Name: "Username", Value: username})
		}
		if instructions != "" {
			fields = append(fields, confidentialField{Name: "Instructions", Value: instructions})
		}
		if secureNotes != "" {
			fields = append(fields, confidentialField{Name: "Confidential Notes", Value: secureNotes})
		}
		items = append(items, confidentialItem{Label: name, Fields: fields})
	}
	return items, nil
}

func (h *ConfidentialHandler) queryServiceAccounts(userID int64) ([]confidentialItem, error) {
	rows, err := h.DB.Query(`
		SELECT name, provider, account_number, secure_notes FROM service_accounts
		WHERE user_id = ? AND (account_number != '' OR secure_notes != '')
		ORDER BY name
	`, userID)
	if err != nil {
		h.Logger.Error("failed to query confidential service accounts", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var items []confidentialItem
	for rows.Next() {
		var name, provider, accountNumber, secureNotes string
		if err := rows.Scan(&name, &provider, &accountNumber, &secureNotes); err != nil {
			continue
		}
		label := name
		if provider != "" {
			label += fmt.Sprintf(" (%s)", provider)
		}
		var fields []confidentialField
		if accountNumber != "" {
			fields = append(fields, confidentialField{Name: "Account Number", Value: accountNumber})
		}
		if secureNotes != "" {
			fields = append(fields, confidentialField{Name: "Confidential Notes", Value: secureNotes})
		}
		items = append(items, confidentialItem{Label: label, Fields: fields})
	}
	return items, nil
}
