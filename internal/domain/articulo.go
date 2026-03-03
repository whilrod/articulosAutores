package domain

import (
	"time"

	"github.com/google/uuid"
)

type Articulo struct {
	ID               uuid.UUID
	Titulo           string
	Contenido        string
	AutorID          uuid.UUID
	Estado           string
	FechaPublicacion time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
