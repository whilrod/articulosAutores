package repositories

import (
	"articulosAutores/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ArticuloRepository struct {
	db *sql.DB
}

func NewArticuloRepository(db *sql.DB) *ArticuloRepository {
	return &ArticuloRepository{db: db}
}

// crea nuevo articulo como borrador
func (r *ArticuloRepository) Create(ctx context.Context, articulo *domain.Articulo) error {
	// Ya no necesitamos validar estado, el dominio lo garantiza
	query := `INSERT INTO articulos (id, titulo, contenido, autor_id, estado, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		articulo.ID.String(),
		articulo.Titulo,
		articulo.Contenido,
		articulo.AutorID.String(),
		articulo.Estado,
		articulo.CreatedAt,
		articulo.UpdatedAt,
	)
	return err
}

// obtiene un artículo a partir de ID
func (a *ArticuloRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Articulo, error) {
	var articulo domain.Articulo
	var idStr, autorIDStr string
	query := `SELECT id, titulo, contenido, autor_id, estado, fecha_publicacion, created_at, updated_at FROM articulos WHERE id = ?`
	row := a.db.QueryRowContext(ctx, query, id.String())
	err := row.Scan(
		&idStr,
		&articulo.Titulo,
		&articulo.Contenido,
		&autorIDStr,
		&articulo.Estado,
		&articulo.FechaPublicacion,
		&articulo.CreatedAt,
		&articulo.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("artículo con ID %s no encontrado: %w", id, err)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener artículo %s: %w", id, err)
	}
	articulo.ID = uuid.MustParse(idStr)
	articulo.AutorID = uuid.MustParse(autorIDStr)
	return &articulo, nil
}

// cambia el estado del artículo a PUBLICADO y asigna fecha
func (a *ArticuloRepository) Publicar(ctx context.Context, id uuid.UUID) (*domain.Articulo, error) {
	articulo, err := a.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := articulo.ValidarParaPublicar(); err != nil {
		return nil, err
	}
	ahora := time.Now()
	articulo.Estado = "publicado"
	articulo.FechaPublicacion = &ahora
	articulo.UpdatedAt = ahora
	query := `UPDATE articulos SET estado = ?, fecha_publicacion = ?, updated_at = ? WHERE id = ?`
	_, err = a.db.ExecContext(ctx, query,
		articulo.Estado,
		articulo.FechaPublicacion,
		articulo.UpdatedAt,
		articulo.ID.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("error al publicar artículo: %w", err)
	}
	return articulo, nil
}

func (a *ArticuloRepository) ListPublicados(ctx context.Context, limit, offset int) ([]*domain.Articulo, int64, error) {
	// Total de artículos publicados
	var total int64
	countQuery := `SELECT COUNT(*) FROM articulos WHERE estado = 'publicado'`
	err := a.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error al contar artículos publicados: %w", err)
	}

	//listado paginado
	query := `SELECT id, titulo, contenido, autor_id, estado, fecha_publicacion, created_at, updated_at FROM articulos estado = 'publicado' ORDER BY fecha_publicacion DESC LIMIT ? OFFSET ?`
	rows, err := a.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error al listar artículos publicados: %w", err)
	}
	defer rows.Close()
	var articulos []*domain.Articulo
	for rows.Next() {
		var articulo domain.Articulo
		var idStr, autorIDStr string
		err := rows.Scan(
			&idStr,
			&articulo.Titulo,
			&articulo.Contenido,
			&autorIDStr,
			&articulo.Estado,
			&articulo.FechaPublicacion,
			&articulo.CreatedAt,
			&articulo.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error al obtener articulo: %w", err)
		}
		articulo.ID = uuid.MustParse(idStr)
		articulo.AutorID = uuid.MustParse(autorIDStr)
		articulos = append(articulos, &articulo)
	}
	return articulos, total, nil
}

// devuelve artículos de un autor, opcionalmente filtrados por estado
func (a *ArticuloRepository) ListByAutor(ctx context.Context, autorID uuid.UUID, estado string, limit, offset int) ([]*domain.Articulo, int64, error) {
	var args []interface{}
	var condiciones string

	if estado != "" {
		condiciones = "WHERE autor_id = ? AND estado = ?"
		args = append(args, autorID.String(), estado)
	} else {
		condiciones = "WHERE autor_id = ?"
		args = append(args, autorID.String())
	}

	// Total de artículos del autor
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM articulos %s", condiciones)
	err := a.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error al contar artículos del autor: %w", err)
	}

	//listado paginado
	query := fmt.Sprintf(`SELECT id, titulo, contenido, autor_id, estado, fecha_publicacion, created_at, updated_at
        FROM articulos %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, condiciones)

	argsConPaginacion := append(args, limit, offset)
	rows, err := a.db.QueryContext(ctx, query, argsConPaginacion...)
	if err != nil {
		return nil, 0, fmt.Errorf("error al listar artículos del autor: %w", err)
	}
	defer rows.Close()
	var articulos []*domain.Articulo
	for rows.Next() {
		var articulo domain.Articulo
		var idStr, autorIDStr string
		err := rows.Scan(
			&idStr,
			&articulo.Titulo,
			&articulo.Contenido,
			&autorIDStr,
			&articulo.Estado,
			&articulo.FechaPublicacion,
			&articulo.CreatedAt,
			&articulo.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error al obtener artículo: %w", err)
		}
		articulo.ID = uuid.MustParse(idStr)
		articulo.AutorID = uuid.MustParse(autorIDStr)
		articulos = append(articulos, &articulo)
	}
	return articulos, total, nil
}

// devuelve estadísticas de artículos de un autor
func (r *ArticuloRepository) GetResumenAutor(ctx context.Context, autorID uuid.UUID) (int, int, *time.Time, error) {
	query := `SELECT COUNT(*) as total, SUM(CASE WHEN estado = 'publicado' THEN 1 ELSE 0 END) as publicados,
            MAX(fecha_publicacion) as ultima FROM articulos WHERE autor_id = ?`

	var total, publicados int
	var ultima *time.Time
	row := r.db.QueryRowContext(ctx, query, autorID.String())
	err := row.Scan(&total, &publicados, &ultima)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("error al obtener resumen del autor %s: %w", autorID, err)
	}
	return total, publicados, ultima, nil
}
