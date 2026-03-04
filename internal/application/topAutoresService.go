package application

import (
	"articulosAutores/internal/domain"
	"articulosAutores/internal/infrastructure/repositories"
	"context"
	"sort"
)

type AutorConScore struct {
	Autor *domain.Autor `json:"autor"`
	Score float64       `json:"score_total"`
}

type TopAutoresService struct {
	autorRepo    *repositories.AutorRepository
	articuloRepo *repositories.ArticuloRepository
}

func NewTopAutoresService(autorRepo *repositories.AutorRepository, articuloRepo *repositories.ArticuloRepository) *TopAutoresService {
	return &TopAutoresService{
		autorRepo:    autorRepo,
		articuloRepo: articuloRepo,
	}
}

// devuelve los autores con mayor score acumulado
func (s *TopAutoresService) GetTop(ctx context.Context, n int) ([]*AutorConScore, error) {
	autores, _, err := s.autorRepo.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}
	var resultados []*AutorConScore
	for _, autor := range autores {
		// Solo necesitamos total de publicados
		_, publicados, _, err := s.articuloRepo.GetResumenAutor(ctx, autor.ID)
		if err != nil {
			continue
		}
		if publicados == 0 {
			resultados = append(resultados, &AutorConScore{
				Autor: autor,
				Score: 0,
			})
			continue
		}

		//limite alto para incluir hasta 10000 articulos por autor y tener score suficiente
		articulos, _, err := s.articuloRepo.ListByAutor(ctx, autor.ID, "publicado", 10000, 0)
		if err != nil {
			continue
		}
		var scoreTotal float64
		for _, art := range articulos {
			scoreTotal += domain.CalcularScore(art, publicados)
		}
		resultados = append(resultados, &AutorConScore{
			Autor: autor,
			Score: scoreTotal,
		})
	}
	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].Score > resultados[j].Score
	})

	if n > len(resultados) {
		n = len(resultados)
	}
	return resultados[:n], nil
}
