package model

import (
	"time"

	"github.com/google/uuid"

	vm "machineIssuerSystem/internal/virtualmachine"
)

type Product struct {
	ID          uuid.UUID
	Title       string
	Description string
	Tags        []string
	ImageURLs   []string
}

type Server struct {
	ID        uuid.UUID
	Title     string
	CPU       int
	Memory    int
	Disk      int
	RentBy    *uuid.UUID
	IP        string
	RentUntil *time.Time
}

type Metric struct {
	Uptime int64
	CPU    float64
	RAM    float64
	Memory int64
}

func FromPkgToDomain(req vm.Metrics) Metric {
	return Metric{
		Uptime: req.Uptime,
		CPU:    req.CPU,
		RAM:    req.RAM,
		Memory: req.MEM,
	}
}

type RentServerRequest struct {
	BookingDays int64
}
