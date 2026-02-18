package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tahiriqbal095/attendify/internal/models"
)

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *models.Enrollment) error
	GetByClassID(ctx context.Context, classID uuid.UUID) ([]models.Enrollment, error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]models.Enrollment, error)
	IsEnrolled(ctx context.Context, classID, studentID uuid.UUID) (bool, error)
	Delete(ctx context.Context, classID, studentID uuid.UUID) error
	GetClassesWithDetailsByStudentID(ctx context.Context, studentID uuid.UUID) ([]models.EnrollmentWithClass, error)
	GetStudentsWithDetailsByClassID(ctx context.Context, classID uuid.UUID) ([]models.StudentInClass, error)
}

type enrollmentRepository struct {
	pool *pgxpool.Pool
}

func NewEnrollmentRepository(pool *pgxpool.Pool) EnrollmentRepository {
	return &enrollmentRepository{pool: pool}
}

func (r *enrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) error {
	query := `
		INSERT INTO enrollments (id, class_id, student_id, enrolled_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.pool.Exec(ctx, query,
		enrollment.ID,
		enrollment.ClassID,
		enrollment.StudentID,
		enrollment.EnrolledAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateKey
		}
		return fmt.Errorf("failed to create enrollment: %w", err)
	}

	return nil
}

func (r *enrollmentRepository) GetByClassID(ctx context.Context, classID uuid.UUID) ([]models.Enrollment, error) {
	query := `
		SELECT id, class_id, student_id, enrolled_at
		FROM enrollments
		WHERE class_id = $1
		ORDER BY enrolled_at DESC
	`

	rows, err := r.pool.Query(ctx, query, classID)
	if err != nil {
		return nil, fmt.Errorf("failed to query enrollments: %w", err)
	}
	defer rows.Close()

	var enrollments []models.Enrollment
	for rows.Next() {
		var e models.Enrollment
		if err := rows.Scan(&e.ID, &e.ClassID, &e.StudentID, &e.EnrolledAt); err != nil {
			return nil, fmt.Errorf("failed to scan enrollment: %w", err)
		}
		enrollments = append(enrollments, e)
	}

	return enrollments, rows.Err()
}

func (r *enrollmentRepository) GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]models.Enrollment, error) {
	query := `
		SELECT id, class_id, student_id, enrolled_at
		FROM enrollments
		WHERE student_id = $1
		ORDER BY enrolled_at DESC
	`

	rows, err := r.pool.Query(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query enrollments: %w", err)
	}
	defer rows.Close()

	var enrollments []models.Enrollment
	for rows.Next() {
		var e models.Enrollment
		if err := rows.Scan(&e.ID, &e.ClassID, &e.StudentID, &e.EnrolledAt); err != nil {
			return nil, fmt.Errorf("failed to scan enrollment: %w", err)
		}
		enrollments = append(enrollments, e)
	}

	return enrollments, rows.Err()
}

func (r *enrollmentRepository) IsEnrolled(ctx context.Context, classID, studentID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM enrollments WHERE class_id = $1 AND student_id = $2)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, classID, studentID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check enrollment: %w", err)
	}

	return exists, nil
}

func (r *enrollmentRepository) Delete(ctx context.Context, classID, studentID uuid.UUID) error {
	query := `DELETE FROM enrollments WHERE class_id = $1 AND student_id = $2`

	result, err := r.pool.Exec(ctx, query, classID, studentID)
	if err != nil {
		return fmt.Errorf("failed to delete enrollment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// GetClassesWithDetailsByStudentID returns enrolled classes with full class details.
func (r *enrollmentRepository) GetClassesWithDetailsByStudentID(ctx context.Context, studentID uuid.UUID) ([]models.EnrollmentWithClass, error) {
	query := `
		SELECT e.id, e.enrolled_at, c.id, c.name, c.code, c.teacher_id, c.created_at
		FROM enrollments e
		JOIN classes c ON e.class_id = c.id
		WHERE e.student_id = $1
		ORDER BY e.enrolled_at DESC
	`

	rows, err := r.pool.Query(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query enrollments with classes: %w", err)
	}
	defer rows.Close()

	var result []models.EnrollmentWithClass
	for rows.Next() {
		var ec models.EnrollmentWithClass
		if err := rows.Scan(
			&ec.ID,
			&ec.EnrolledAt,
			&ec.Class.ID,
			&ec.Class.Name,
			&ec.Class.Code,
			&ec.Class.TeacherID,
			&ec.Class.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan enrollment with class: %w", err)
		}
		result = append(result, ec)
	}

	return result, rows.Err()
}

// GetStudentsWithDetailsByClassID returns enrolled students with full user details.
func (r *enrollmentRepository) GetStudentsWithDetailsByClassID(
	ctx context.Context, classID uuid.UUID,
) ([]models.StudentInClass, error) {
	query := `
		SELECT e.id, e.enrolled_at, u.id, u.email, u.name, u.role, u.created_at
		FROM enrollments e
		JOIN users u ON e.student_id = u.id
		WHERE e.class_id = $1
		ORDER BY u.name ASC
	`

	rows, err := r.pool.Query(ctx, query, classID)
	if err != nil {
		return nil, fmt.Errorf("failed to query students in class: %w", err)
	}
	defer rows.Close()

	var result []models.StudentInClass
	for rows.Next() {
		var sc models.StudentInClass
		if err := rows.Scan(
			&sc.ID,
			&sc.EnrolledAt,
			&sc.Student.ID,
			&sc.Student.Email,
			&sc.Student.Name,
			&sc.Student.Role,
			&sc.Student.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan student in class: %w", err)
		}
		result = append(result, sc)
	}

	return result, rows.Err()
}
