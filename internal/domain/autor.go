package domain

import (
	"time"

	"github.com/google/uuid"
)

type Autor struct {
	ID            uuid.UUID
	Nombre        string
	Email         string
	Bio           string
	FechaRegistro time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
