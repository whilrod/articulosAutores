package handlers

import (
	"net/http"
	"strconv"

	"articulosAutores/internal/application"
	"articulosAutores/internal/infrastructure/handlers/dto"

	"github.com/gin-gonic/gin"
)

type TopAutoresHandler struct {
	service *application.TopAutoresService
}

func NewTopAutoresHandler(service *application.TopAutoresService) *TopAutoresHandler {
	return &TopAutoresHandler{service: service}
}

// GetTopAutores godoc
// @Summary Obtener los N autores con mayor score acumulado
// @Tags autores
// @Produce json
// @Param n query int false "Número de autores a retornar" default(3)
// @Success 200 {object} dto.TopAutoresResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/autores/top [get]
func (h *TopAutoresHandler) GetTopAutores(c *gin.Context) {
	n, _ := strconv.Atoi(c.DefaultQuery("n", "3"))

	if n < 1 {
		n = 3
	}
	resultados, err := h.service.GetTop(c.Request.Context(), n)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al obtener top autores",
		})
		return
	}
	data := make([]*dto.TopAutorResponse, len(resultados))
	for i, r := range resultados {
		data[i] = &dto.TopAutorResponse{
			Autor: &dto.AutorResponse{
				ID:            r.Autor.ID,
				Nombre:        r.Autor.Nombre,
				Email:         r.Autor.Email,
				Bio:           r.Autor.Bio,
				FechaRegistro: r.Autor.FechaRegistro,
				CreatedAt:     r.Autor.CreatedAt,
				UpdatedAt:     r.Autor.UpdatedAt,
			},
			Score: r.Score,
		}
	}
	response := dto.TopAutoresResponse{
		Data: data,
	}
	c.JSON(http.StatusOK, response)
}
