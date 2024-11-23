package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID         uuid.UUID `db:"id" json:"uuid"`
	Email        string
	Username     string
	HashPassword string `db:"password"`
	Role         int64
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
