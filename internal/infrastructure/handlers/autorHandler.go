package handlers

import (
	"net/http"

	"articulosAutores/internal/domain"
	"articulosAutores/internal/infrastructure/handlers/dto"
	"articulosAutores/internal/infrastructure/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AutorHandler struct {
	repo *repositories.AutorRepository
}

func NewAutorHandler(repo *repositories.AutorRepository) *AutorHandler {
	return &AutorHandler{repo: repo}
}

// CreateAutor godoc
// @Summary Crear un nuevo autor
// @Tags autores
// @Accept json
// @Produce json
// @Param request body dto.CrearAutorRequest true "Datos del autor"
// @Success 201 {object} dto.AutorResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/autores [post]
func (h *AutorHandler) CreateAutor(c *gin.Context) {
	var req dto.CrearAutorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"detalle": err.Error(),
		})
		return
	}

	// Crear autor desde dominio
	autor := domain.NewAutor(req.Nombre, req.Email, req.Bio)
	if err := h.repo.Create(c.Request.Context(), autor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al crear autor",
		})
		return
	}

	// Convertir a DTO de respuesta
	response := dto.AutorResponse{
		ID:            autor.ID,
		Nombre:        autor.Nombre,
		Email:         autor.Email,
		Bio:           autor.Bio,
		FechaRegistro: autor.FechaRegistro,
		CreatedAt:     autor.CreatedAt,
		UpdatedAt:     autor.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetAutorByID godoc
// @Summary Obtener autor por ID
// @Tags autores
// @Produce json
// @Param id path string true "ID del autor"
// @Success 200 {object} dto.AutorResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/autores/{id} [get]
func (h *AutorHandler) GetAutorByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}
	autor, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Autor no encontrado",
		})
		return
	}
	response := dto.AutorResponse{
		ID:            autor.ID,
		Nombre:        autor.Nombre,
		Email:         autor.Email,
		Bio:           autor.Bio,
		FechaRegistro: autor.FechaRegistro,
		CreatedAt:     autor.CreatedAt,
		UpdatedAt:     autor.UpdatedAt,
	}
	c.JSON(http.StatusOK, response)
}
