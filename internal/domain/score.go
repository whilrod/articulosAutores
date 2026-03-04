package domain

import "time"

//calcula el score de relevancia de un artículo
func CalcularScore(articulo *Articulo, totalPublicadosAutor int) float64 {
	if articulo.Estado != "publicado" || articulo.FechaPublicacion == nil {
		return 0
	}
	palabras := float64(articulo.ContarPalabras())
	score := palabras * 0.1
	score += float64(totalPublicadosAutor) * 5
	score += calcularBonusReciente(*articulo.FechaPublicacion)
	return score
}

//determina el bonus según tiempo desde publicación
func calcularBonusReciente(fechaPublicacion time.Time) float64 {
	horas := time.Since(fechaPublicacion).Hours()
	switch {
	case horas < 24:
		return 50
	case horas < 72:
		return 20
	default:
		return 0
	}
}
