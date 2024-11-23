package model

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID
	Title       string
	Description string
	Tags        []string
	ImageURLs   []string
}

type Server struct {
	ID     uuid.UUID
	Title  string
	CPU    int
	Memory int
	Disk   int
	RentBy *uuid.UUID
	IP     string
}
