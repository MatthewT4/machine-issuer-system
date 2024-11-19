package storage

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID
	Title       string
	Description string
	Tags        []string
	ImageURL    []string `db:"image_urls"`
}
