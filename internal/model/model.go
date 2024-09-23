package model

import (
	"errors"
	"time"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrIdentifierExists = errors.New("identifier already exists")
)

type Shortening struct {
	Identifier  string    `gorm:"primaryKey" json:"identifier"`
	OriginalURL string    `gorm:"not null" json:"original_url"`
	Visits      int64     `gorm:"default:0" json:"visits"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ShortenInput struct {
	RawURL     string `json:"raw_url"`
	Identifier string `json:"identifier,omitempty"`
	CreateBy   string `json:"create_by"`
}
