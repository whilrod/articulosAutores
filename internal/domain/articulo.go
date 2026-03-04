package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Articulo struct {
	ID               uuid.UUID
	Titulo           string
	Contenido        string
	AutorID          uuid.UUID
	Estado           string
	FechaPublicacion *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewArticuloBorrador crea un nuevo artículo en estado BORRADOR
func NewArticuloBorrador(titulo, contenido string, autorID uuid.UUID) *Articulo {
	ahora := time.Now()

	return &Articulo{
		ID:        uuid.New(),
		Titulo:    titulo,
		Contenido: contenido,
		AutorID:   autorID,
		Estado:    "borrador",
		CreatedAt: ahora,
		UpdatedAt: ahora,
	}
}

func (a *Articulo) ContarPalabras() int {
	palabras := strings.Fields(a.Contenido)
	return len(palabras)
}

func (a *Articulo) CalcularRepetidas() float64 {
	palabras := strings.Fields(a.Contenido)
	if len(palabras) == 0 {
		return 0
	}

	// Contar frecuencia de palabras
	frecuencia := make(map[string]int)
	for _, p := range palabras {
		frecuencia[p]++
	}

	// Contar palabras que aparecen más de una vez
	repetidas := 0
	for _, count := range frecuencia {
		if count > 1 {
			repetidas += count
		}
	}
	return float64(repetidas) / float64(len(palabras)) * 100
}

func (a *Articulo) ValidarParaPublicar() error {
	if a.Estado == "publicado" {
		return fmt.Errorf("el artículo ya está publicado")
	}
	palabras := a.ContarPalabras()
	if palabras < 120 {
		return fmt.Errorf("el artículo debe tener mínimo 120 palabras (tiene %d)", palabras)
	}
	repetidas := a.CalcularRepetidas()
	if repetidas > 35 {
		return fmt.Errorf("el artículo supera el 35%% de palabras repetidas (tiene %.1f%%)", repetidas)
	}
	return nil
}

func (a *Articulo) Publicar() {
	ahora := time.Now()
	a.Estado = "publicado"
	a.FechaPublicacion = &ahora
	a.UpdatedAt = ahora
}
