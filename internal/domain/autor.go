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

// crea un nuevo autor
func NewAutor(nombre, email, bio string) *Autor {
	ahora := time.Now()

	return &Autor{
		ID:            uuid.New(),
		Nombre:        nombre,
		Email:         email,
		Bio:           bio,
		FechaRegistro: ahora,
		CreatedAt:     ahora,
		UpdatedAt:     ahora,
	}
}
