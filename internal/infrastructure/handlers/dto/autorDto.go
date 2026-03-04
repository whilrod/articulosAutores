package dto

import (
	"time"

	"github.com/google/uuid"
)

type CrearAutorRequest struct {
	Nombre string `json:"nombre" binding:"required,min=2,max=100"`
	Email  string `json:"email" binding:"required,email"`
	Bio    string `json:"bio" binding:"max=500"`
}

type AutorResponse struct {
	ID            uuid.UUID `json:"id"`
	Nombre        string    `json:"nombre"`
	Email         string    `json:"email"`
	Bio           string    `json:"bio,omitempty"`
	FechaRegistro time.Time `json:"fecha_registro"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
