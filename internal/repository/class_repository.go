package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tahiriqbal095/attendify/internal/models"
)

type ClassRepository interface {
	Create(ctx context.Context, class *models.Class) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Class, error)
	GetByCode(ctx context.Context, code string) (*models.Class, error)
	GetByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]models.Class, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type classRepository struct {
	pool *pgxpool.Pool
}

func NewClassRepository(pool *pgxpool.Pool) ClassRepository {
	return &classRepository{pool: pool}
}

func (r *classRepository) Create(ctx context.Context, class *models.Class) error {
	query := `
		INSERT INTO classes (id, name, code, teacher_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, query,
		class.ID,
		class.Name,
		class.Code,
		class.TeacherID,
		class.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateKey
		}
		return fmt.Errorf("failed to create class: %w", err)
	}

	return nil
}

func (r *classRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Class, error) {
	query := `
		SELECT id, name, code, teacher_id, created_at
		FROM classes
		WHERE id = $1
	`

	class := &models.Class{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&class.ID,
		&class.Name,
		&class.Code,
		&class.TeacherID,
		&class.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get class by id: %w", err)
	}

	return class, nil
}

func (r *classRepository) GetByCode(ctx context.Context, code string) (*models.Class, error) {
	query := `
		SELECT id, name, code, teacher_id, created_at
		FROM classes
		WHERE code = $1
	`

	class := &models.Class{}
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&class.ID,
		&class.Name,
		&class.Code,
		&class.TeacherID,
		&class.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get class by code: %w", err)
	}

	return class, nil
}

func (r *classRepository) GetByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]models.Class, error) {
	query := `
		SELECT id, name, code, teacher_id, created_at
		FROM classes
		WHERE teacher_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to query classes: %w", err)
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var class models.Class
		if err := rows.Scan(
			&class.ID,
			&class.Name,
			&class.Code,
			&class.TeacherID,
			&class.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan class: %w", err)
		}
		classes = append(classes, class)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating classes: %w", err)
	}

	return classes, nil
}

func (r *classRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM classes WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete class: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
