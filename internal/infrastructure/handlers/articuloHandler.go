package handlers

import (
	"net/http"
	"strconv"

	"articulosAutores/internal/domain"
	"articulosAutores/internal/infrastructure/handlers/dto"
	"articulosAutores/internal/infrastructure/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ArticuloHandler struct {
	repo *repositories.ArticuloRepository
}

func NewArticuloHandler(repo *repositories.ArticuloRepository) *ArticuloHandler {
	return &ArticuloHandler{repo: repo}
}

// CreateArticulo godoc
// @Summary Crear un nuevo artículo en estado BORRADOR
// @Tags articulos
// @Accept json
// @Produce json
// @Param request body dto.CrearArticuloRequest true "Datos del artículo"
// @Success 201 {object} dto.ArticuloResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/articulos [post]
func (h *ArticuloHandler) CreateArticulo(c *gin.Context) {
	var req dto.CrearArticuloRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"detalle": err.Error(),
		})
		return
	}
	articulo := domain.NewArticuloBorrador(req.Titulo, req.Contenido, req.AutorID)
	if err := h.repo.Create(c.Request.Context(), articulo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al crear artículo",
		})
		return
	}
	response := dto.ArticuloResponse{
		ID:        articulo.ID,
		Titulo:    articulo.Titulo,
		Contenido: articulo.Contenido,
		AutorID:   articulo.AutorID,
		Estado:    articulo.Estado,
		CreatedAt: articulo.CreatedAt,
		UpdatedAt: articulo.UpdatedAt,
	}
	c.JSON(http.StatusCreated, response)
}

// PublicarArticulo godoc
// @Summary Publicar un artículo (cambia estado a PUBLICADO)
// @Tags articulos
// @Produce json
// @Param id path string true "ID del artículo"
// @Success 200 {object} dto.ArticuloResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /api/v1/articulos/{id}/publicar [post]
func (h *ArticuloHandler) PublicarArticulo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}
	articulo, err := h.repo.Publicar(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "el artículo debe tener mínimo 120 palabras" ||
			err.Error() == "el artículo supera el 35% de palabras repetidas" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Artículo no encontrado",
		})
		return
	}
	response := dto.ArticuloResponse{
		ID:               articulo.ID,
		Titulo:           articulo.Titulo,
		Contenido:        articulo.Contenido,
		AutorID:          articulo.AutorID,
		Estado:           articulo.Estado,
		FechaPublicacion: articulo.FechaPublicacion,
		CreatedAt:        articulo.CreatedAt,
		UpdatedAt:        articulo.UpdatedAt,
	}
	c.JSON(http.StatusOK, response)
}

// ListArticulosPublicados godoc
// @Summary Listar artículos PUBLICADOS con paginación
// @Tags articulos
// @Produce json
// @Param pagina query int false "Número de página" default(1)
// @Param limite query int false "Elementos por página" default(10)
// @Success 200 {object} dto.ListaArticulosResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/articulos [get]
func (h *ArticuloHandler) ListArticulosPublicados(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	limite, _ := strconv.Atoi(c.DefaultQuery("limite", "10"))
	if pagina < 1 {
		pagina = 1
	}
	if limite < 1 || limite > 100 {
		limite = 10
	}
	offset := (pagina - 1) * limite
	articulos, total, err := h.repo.ListPublicados(c.Request.Context(), limite, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al listar artículos",
		})
		return
	}
	data := make([]*dto.ArticuloResponse, len(articulos))
	for i, art := range articulos {
		data[i] = &dto.ArticuloResponse{
			ID:               art.ID,
			Titulo:           art.Titulo,
			Contenido:        art.Contenido,
			AutorID:          art.AutorID,
			Estado:           art.Estado,
			FechaPublicacion: art.FechaPublicacion,
			CreatedAt:        art.CreatedAt,
			UpdatedAt:        art.UpdatedAt,
		}
	}
	totalPaginas := int(total) / limite
	if int(total)%limite > 0 {
		totalPaginas++
	}
	response := dto.ListaArticulosResponse{
		Data:         data,
		Total:        total,
		Pagina:       pagina,
		Limite:       limite,
		TotalPaginas: totalPaginas,
	}
	c.JSON(http.StatusOK, response)
}

// ListArticulosByAutor godoc
// @Summary Listar artículos por autor (opcional filtrar por estado)
// @Tags articulos
// @Produce json
// @Param id path string true "ID del autor"
// @Param estado query string false "Filtrar por estado (borrador, publicado, archivado)"
// @Param pagina query int false "Número de página" default(1)
// @Param limite query int false "Elementos por página" default(10)
// @Success 200 {object} dto.ListaArticulosResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/autores/{id}/articulos [get]
func (h *ArticuloHandler) ListArticulosByAutor(c *gin.Context) {
	autorIDStr := c.Param("id")
	autorID, err := uuid.Parse(autorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de autor inválido",
		})
		return
	}
	estado := c.Query("estado")
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	limite, _ := strconv.Atoi(c.DefaultQuery("limite", "10"))
	if pagina < 1 {
		pagina = 1
	}
	if limite < 1 || limite > 100 {
		limite = 10
	}
	offset := (pagina - 1) * limite
	articulos, total, err := h.repo.ListByAutor(c.Request.Context(), autorID, estado, limite, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al listar artículos",
		})
		return
	}
	data := make([]*dto.ArticuloResponse, len(articulos))
	for i, art := range articulos {
		data[i] = &dto.ArticuloResponse{
			ID:               art.ID,
			Titulo:           art.Titulo,
			Contenido:        art.Contenido,
			AutorID:          art.AutorID,
			Estado:           art.Estado,
			FechaPublicacion: art.FechaPublicacion,
			CreatedAt:        art.CreatedAt,
			UpdatedAt:        art.UpdatedAt,
		}
	}
	totalPaginas := int(total) / limite
	if int(total)%limite > 0 {
		totalPaginas++
	}
	response := dto.ListaArticulosResponse{
		Data:         data,
		Total:        total,
		Pagina:       pagina,
		Limite:       limite,
		TotalPaginas: totalPaginas,
	}
	c.JSON(http.StatusOK, response)
}

// GetResumenAutor godoc
// @Summary Obtener resumen de artículos de un autor
// @Tags autores
// @Produce json
// @Param id path string true "ID del autor"
// @Success 200 {object} dto.ResumenAutorResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/autores/{id}/resumen [get]
func (h *ArticuloHandler) GetResumenAutor(c *gin.Context) {
	autorIDStr := c.Param("id")

	autorID, err := uuid.Parse(autorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de autor inválido",
		})
		return
	}
	total, publicados, ultima, err := h.repo.GetResumenAutor(c.Request.Context(), autorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al obtener resumen",
		})
		return
	}
	response := dto.ResumenAutorResponse{
		TotalArticulos:    total,
		TotalPublicados:   publicados,
		UltimaPublicacion: ultima,
	}
	c.JSON(http.StatusOK, response)
}
