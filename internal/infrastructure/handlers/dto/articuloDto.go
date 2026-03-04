package dto

import (
	"time"

	"github.com/google/uuid"
)

type CrearArticuloRequest struct {
	Titulo    string    `json:"titulo" binding:"required,min=5,max=200"`
	Contenido string    `json:"contenido" binding:"required,min=50"`
	AutorID   uuid.UUID `json:"autor_id" binding:"required"`
}

type PublicarArticuloResponse struct {
	Mensaje string `json:"mensaje"`
}

type ArticuloResponse struct {
	ID               uuid.UUID  `json:"id"`
	Titulo           string     `json:"titulo"`
	Contenido        string     `json:"contenido"`
	AutorID          uuid.UUID  `json:"autor_id"`
	Estado           string     `json:"estado"`
	FechaPublicacion *time.Time `json:"fecha_publicacion,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Score            *float64   `json:"score,omitempty"`
}

type ListaArticulosResponse struct {
	Data         []*ArticuloResponse `json:"data"`
	Total        int64               `json:"total"`
	Pagina       int                 `json:"pagina"`
	Limite       int                 `json:"limite"`
	TotalPaginas int                 `json:"total_paginas"`
}

type ResumenAutorResponse struct {
	TotalArticulos    int        `json:"total_articulos"`
	TotalPublicados   int        `json:"total_publicados"`
	UltimaPublicacion *time.Time `json:"ultima_publicacion,omitempty"`
}

type TopAutorResponse struct {
	Autor *AutorResponse `json:"autor"`
	Score float64        `json:"score_total"`
}

type TopAutoresResponse struct {
	Data []*TopAutorResponse `json:"data"`
}
