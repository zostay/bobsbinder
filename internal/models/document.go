package models

import "time"

type DocumentCategory struct {
	ID          int64  `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	SortOrder   int    `json:"sort_order"`
}

type Document struct {
	ID         int64     `json:"id"`
	PartyID    int64     `json:"party_id"`
	CategoryID int64     `json:"category_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DocumentFile struct {
	ID          int64     `json:"id"`
	DocumentID  int64     `json:"document_id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	FilePath    string    `json:"-"`
	FileSize    int64     `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
}

type ChecklistItem struct {
	ID         int64     `json:"id"`
	PartyID    int64     `json:"party_id"`
	CategoryID int64     `json:"category_id"`
	Completed  bool      `json:"completed"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Share struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Token          string    `json:"token"`
	UnlockCodeHash string    `json:"-"`
	ExpiresAt      time.Time `json:"expires_at"`
	CreatedAt      time.Time `json:"created_at"`
}
