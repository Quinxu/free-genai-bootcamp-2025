package models

import (
	"time"
)

// Base contains common fields for all models
type Base struct {
	ID        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
