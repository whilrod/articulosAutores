package repositories

import (
	"articulosAutores/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AutorRepository struct {
	db *sql.DB
}

func NewAutorRepository(db *sql.DB) *AutorRepository {
	return &AutorRepository{db: db}
}

// crear nuevo autor en db
func (a *AutorRepository) Create(ctx context.Context, autor *domain.Autor) error {
	query := `INSTERT INTO autores (id,nombre,email,bio,fecha_registro,created_at,updated_at) VALUES (?,?,?,?,?,?,?)`

	_, err := a.db.ExecContext(
		ctx, query,
		autor.ID.String(),
		autor.Nombre,
		autor.Email,
		autor.Bio,
		autor.FechaRegistro,
		autor.CreatedAt,
		autor.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creando autor %w", err)
	}
	return nil
}

// obtener autor a partir de id
func (a *AutorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Autor, error) {
	var autor domain.Autor
	var bio sql.NullString
	var idStr string
	query := `SELECT id, nombre, email, bio, fecha_registro, created_at, updated_at FROM autores WHERE id = ?`

	err := a.db.QueryRowContext(
		ctx, query,
		id.String()).Scan(
		&idStr,
		&autor.Nombre,
		&autor.Email,
		&bio,
		&autor.FechaRegistro,
		&autor.CreatedAt,
		&autor.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("error autor id: %s no encontrado: %w", id, err)
	}
	if err != nil {
		return nil, fmt.Errorf("error obteniendo autor por ID: %w", err)
	}
	autor.ID = uuid.MustParse(idStr)
	if bio.Valid {
		autor.Bio = bio.String
	}
	return &autor, nil
}

// 0btener autor a partir de email
func (a *AutorRepository) GetByEmail(ctx context.Context, email string) (*domain.Autor, error) {
	var autor domain.Autor
	var bio sql.NullString
	var idStr string

	query := `SELECT id, nombre, email, bio, fecha_registro, created_at, updated_at FROM autores WHERE email = ?`
	row := a.db.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&idStr,
		&autor.Nombre,
		&autor.Email,
		&bio,
		&autor.FechaRegistro,
		&autor.CreatedAt,
		&autor.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("autor con email %s no encontrado: %w", email, err)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener autor por email %s: %w", email, err)
	}

	autor.ID = uuid.MustParse(idStr)
	if bio.Valid {
		autor.Bio = bio.String
	}

	return &autor, nil
}

func (a *AutorRepository) Update(ctx context.Context, autor *domain.Autor) error {
	autor.UpdatedAt = time.Now()

	query := `UPDATE autores SET nombre = ?, email = ?, bio = ?, updated_at = ? WHERE id = ?`
	result, err := a.db.ExecContext(
		ctx, query,
		autor.Nombre,
		autor.Email,
		autor.Bio,
		autor.UpdatedAt,
		autor.ID.String(),
	)
	if err != nil {
		return fmt.Errorf("error al actualizar autor %s: %w", autor.ID, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas para autor %s: %w", autor.ID, err)
	}
	if rows == 0 {
		return fmt.Errorf("autor con ID %s no encontrado para actualizar", autor.ID)
	}
	return nil
}

// eliminar autor a partir de su id
func (a *AutorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM autores WHERE id = ?`
	result, err := a.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("error al eliminar autor %s: %w", id, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas al eliminar autor %s: %w", id, err)
	}
	if rows == 0 {
		return fmt.Errorf("autor con ID %s no encontrado para eliminar", id)
	}
	return nil
}
