package models

import "time"

type Contact struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Name         string    `json:"name"`
	Relationship string    `json:"relationship,omitempty"`
	Role         string    `json:"role,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Email        string    `json:"email,omitempty"`
	Address      string    `json:"address,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Location struct {
	ID                  int64     `json:"id"`
	UserID              int64     `json:"user_id"`
	Name                string    `json:"name"`
	Type                string    `json:"type"`
	Description         string    `json:"description,omitempty"`
	Address             string    `json:"address,omitempty"`
	AccessInstructions  string    `json:"access_instructions,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type DigitalAccess struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Username     string    `json:"username,omitempty"`
	Instructions string    `json:"instructions,omitempty"`
	LocationID   *int64    `json:"location_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type InsurancePolicy struct {
	ID             int64    `json:"id"`
	UserID         int64    `json:"user_id"`
	PartyID        *int64   `json:"party_id,omitempty"`
	Provider       string   `json:"provider"`
	PolicyNumber   string   `json:"policy_number,omitempty"`
	Type           string   `json:"type,omitempty"`
	CoverageAmount *float64 `json:"coverage_amount,omitempty"`
	Beneficiary    string   `json:"beneficiary,omitempty"`
	AgentName      string   `json:"agent_name,omitempty"`
	AgentPhone     string   `json:"agent_phone,omitempty"`
	LocationID     *int64   `json:"location_id,omitempty"`
	Notes          string   `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ServiceAccount struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Type          string    `json:"type"`
	Name          string    `json:"name"`
	Provider      string    `json:"provider,omitempty"`
	AccountNumber string    `json:"account_number,omitempty"`
	ContactName   string    `json:"contact_name,omitempty"`
	ContactPhone  string    `json:"contact_phone,omitempty"`
	ContactEmail  string    `json:"contact_email,omitempty"`
	Notes         string    `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PartyObituaryInfo struct {
	ID           int64     `json:"id"`
	PartyID      int64     `json:"party_id"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Relationship string    `json:"relationship,omitempty"`
	Details      string    `json:"details,omitempty"`
	EventDate    *string   `json:"event_date,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SurvivorLetter struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Greeting  string    `json:"greeting"`
	Intro     string    `json:"intro"`
	Closing   string    `json:"closing"`
	Signature string    `json:"signature"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SurvivorLetterSection struct {
	ID         int64                `json:"id"`
	LetterID   int64                `json:"letter_id"`
	SectionKey string               `json:"section_key"`
	Title      string               `json:"title"`
	SortOrder  int                  `json:"sort_order"`
	Visible    bool                 `json:"visible"`
	Items      []SurvivorLetterItem `json:"items"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

type SurvivorLetterItem struct {
	ID         int64     `json:"id"`
	SectionID  int64     `json:"section_id"`
	SourceType *string   `json:"source_type,omitempty"`
	SourceID   *int64    `json:"source_id,omitempty"`
	Content    string    `json:"content"`
	ItemType   string    `json:"item_type"`
	Provenance string    `json:"provenance"`
	Suppressed bool      `json:"suppressed"`
	SortOrder  int       `json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FullSurvivorLetter struct {
	SurvivorLetter
	Sections []SurvivorLetterSection `json:"sections"`
}
