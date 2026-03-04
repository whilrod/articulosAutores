package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"articulosAutores/internal/application"
	"articulosAutores/internal/infrastructure/database"
	"articulosAutores/internal/infrastructure/handlers"
	"articulosAutores/internal/infrastructure/repositories"
)

func main() {
	cfg := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	// Conectar a MySQL
	db, err := database.NewConnectionMySql(cfg)
	if err != nil {
		log.Fatalf("Error conectando a base de datos: %v", err)
	}
	defer db.Close()
	autorRepo := repositories.NewAutorRepository(db)
	articuloRepo := repositories.NewArticuloRepository(db)
	topAutoresService := application.NewTopAutoresService(autorRepo, articuloRepo)
	autorHandler := handlers.NewAutorHandler(autorRepo)
	articuloHandler := handlers.NewArticuloHandler(articuloRepo)
	topAutoresHandler := handlers.NewTopAutoresHandler(topAutoresService)
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "error", "database": "disconnected"})
			return
		}
		c.JSON(200, gin.H{"status": "ok", "database": "connected"})
	})
	v1 := router.Group("/api/v1")
	{
		autores := v1.Group("/autores")
		{
			autores.POST("/", autorHandler.CreateAutor)
			autores.GET("/:id", autorHandler.GetAutorByID)
			autores.GET("/:id/articulos", articuloHandler.ListArticulosByAutor)
			autores.GET("/:id/resumen", articuloHandler.GetResumenAutor)
			autores.GET("/top", topAutoresHandler.GetTopAutores)
		}
		articulos := v1.Group("/articulos")
		{
			articulos.POST("/", articuloHandler.CreateArticulo)
			articulos.GET("/", articuloHandler.ListArticulosPublicados)
			articulos.POST("/:id/publicar", articuloHandler.PublicarArticulo)
		}
	}

	// Iniciar servidor
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor iniciado en puerto %s", port)
	router.Run(":" + port)
}
