package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type InsurancePolicyHandler struct {
	DB     *sql.DB
	Logger *zap.Logger
}

type insurancePolicyRequest struct {
	PartyID        *int64   `json:"party_id"`
	Provider       string   `json:"provider"`
	PolicyNumber   string   `json:"policy_number"`
	Type           string   `json:"type"`
	CoverageAmount *float64 `json:"coverage_amount"`
	Beneficiary    string   `json:"beneficiary"`
	AgentName      string   `json:"agent_name"`
	AgentPhone     string   `json:"agent_phone"`
	LocationID     *int64   `json:"location_id"`
	Notes          string   `json:"notes"`
}

func (h *InsurancePolicyHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(`
		SELECT id, user_id, party_id, provider, policy_number, type, coverage_amount,
		       beneficiary, agent_name, agent_phone, location_id, notes, created_at, updated_at
		FROM insurance_policies WHERE user_id = ? ORDER BY provider
	`, userID)
	if err != nil {
		h.Logger.Error("failed to list insurance policies", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var policies []map[string]any
	for rows.Next() {
		var id, uid int64
		var partyID, locationID sql.NullInt64
		var provider, policyNumber, pType, beneficiary, agentName, agentPhone, notes, createdAt, updatedAt string
		var coverageAmount sql.NullFloat64
		if err := rows.Scan(&id, &uid, &partyID, &provider, &policyNumber, &pType, &coverageAmount,
			&beneficiary, &agentName, &agentPhone, &locationID, &notes, &createdAt, &updatedAt); err != nil {
			h.Logger.Error("failed to scan insurance policy", zap.Error(err))
			continue
		}
		item := map[string]any{
			"id": id, "user_id": uid, "provider": provider, "policy_number": policyNumber,
			"type": pType, "beneficiary": beneficiary, "agent_name": agentName,
			"agent_phone": agentPhone, "notes": notes,
			"created_at": createdAt, "updated_at": updatedAt,
		}
		if partyID.Valid {
			item["party_id"] = partyID.Int64
		} else {
			item["party_id"] = nil
		}
		if coverageAmount.Valid {
			item["coverage_amount"] = coverageAmount.Float64
		} else {
			item["coverage_amount"] = nil
		}
		if locationID.Valid {
			item["location_id"] = locationID.Int64
		} else {
			item["location_id"] = nil
		}
		policies = append(policies, item)
	}

	if policies == nil {
		policies = []map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policies)
}

func (h *InsurancePolicyHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	policyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid policy id"}`, http.StatusBadRequest)
		return
	}

	var id, uid int64
	var partyID, locationID sql.NullInt64
	var provider, policyNumber, pType, beneficiary, agentName, agentPhone, notes, createdAt, updatedAt string
	var coverageAmount sql.NullFloat64
	err = h.DB.QueryRow(`
		SELECT id, user_id, party_id, provider, policy_number, type, coverage_amount,
		       beneficiary, agent_name, agent_phone, location_id, notes, created_at, updated_at
		FROM insurance_policies WHERE id = ? AND user_id = ?
	`, policyID, userID).Scan(&id, &uid, &partyID, &provider, &policyNumber, &pType, &coverageAmount,
		&beneficiary, &agentName, &agentPhone, &locationID, &notes, &createdAt, &updatedAt)
	if err != nil {
		http.Error(w, `{"error":"insurance policy not found"}`, http.StatusNotFound)
		return
	}

	item := map[string]any{
		"id": id, "user_id": uid, "provider": provider, "policy_number": policyNumber,
		"type": pType, "beneficiary": beneficiary, "agent_name": agentName,
		"agent_phone": agentPhone, "notes": notes,
		"created_at": createdAt, "updated_at": updatedAt,
	}
	if partyID.Valid {
		item["party_id"] = partyID.Int64
	} else {
		item["party_id"] = nil
	}
	if coverageAmount.Valid {
		item["coverage_amount"] = coverageAmount.Float64
	} else {
		item["coverage_amount"] = nil
	}
	if locationID.Valid {
		item["location_id"] = locationID.Int64
	} else {
		item["location_id"] = nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *InsurancePolicyHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req insurancePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		INSERT INTO insurance_policies (user_id, party_id, provider, policy_number, type, coverage_amount,
		       beneficiary, agent_name, agent_phone, location_id, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.PartyID, req.Provider, req.PolicyNumber, req.Type, req.CoverageAmount,
		req.Beneficiary, req.AgentName, req.AgentPhone, req.LocationID, req.Notes,
	)
	if err != nil {
		h.Logger.Error("failed to create insurance policy", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id": id, "user_id": userID, "party_id": req.PartyID, "provider": req.Provider,
		"policy_number": req.PolicyNumber, "type": req.Type, "coverage_amount": req.CoverageAmount,
		"beneficiary": req.Beneficiary, "agent_name": req.AgentName, "agent_phone": req.AgentPhone,
		"location_id": req.LocationID, "notes": req.Notes,
	})
}

func (h *InsurancePolicyHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	policyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid policy id"}`, http.StatusBadRequest)
		return
	}

	var req insurancePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE insurance_policies SET party_id = ?, provider = ?, policy_number = ?, type = ?,
		       coverage_amount = ?, beneficiary = ?, agent_name = ?, agent_phone = ?,
		       location_id = ?, notes = ?
		WHERE id = ? AND user_id = ?
	`, req.PartyID, req.Provider, req.PolicyNumber, req.Type, req.CoverageAmount,
		req.Beneficiary, req.AgentName, req.AgentPhone, req.LocationID, req.Notes, policyID, userID)
	if err != nil {
		h.Logger.Error("failed to update insurance policy", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"insurance policy not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *InsurancePolicyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	policyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid policy id"}`, http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec("DELETE FROM insurance_policies WHERE id = ? AND user_id = ?", policyID, userID)
	if err != nil {
		h.Logger.Error("failed to delete insurance policy", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"insurance policy not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
