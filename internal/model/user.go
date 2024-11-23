package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID         uuid.UUID
	Email        string
	Username     string
	HashPassword string
	Role         int64
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
