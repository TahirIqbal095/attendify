package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tahiriqbal095/attendify/internal/models"
	"github.com/tahiriqbal095/attendify/internal/repository"
)

type EnrollmentService interface {
	EnrollByCode(ctx context.Context, classCode string, studentID uuid.UUID) (*models.Enrollment, error)
	GetStudentClasses(ctx context.Context, studentID uuid.UUID) ([]models.EnrollmentWithClass, error)
	GetClassStudents(ctx context.Context, classID uuid.UUID) ([]models.StudentInClass, error)
	Unenroll(ctx context.Context, classID, studentID uuid.UUID) error
	IsEnrolled(ctx context.Context, classID, studentID uuid.UUID) (bool, error)
}

type enrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
	classRepo      repository.ClassRepository
}

func NewEnrollmentService(
	enrollmentRepo repository.EnrollmentRepository,
	classRepo repository.ClassRepository,
) EnrollmentService {
	return &enrollmentService{
		enrollmentRepo: enrollmentRepo,
		classRepo:      classRepo,
	}
}

// EnrollByCode enrolls a student in a class using the class code.
func (s *enrollmentService) EnrollByCode(
	ctx context.Context, classCode string, studentID uuid.UUID,
) (*models.Enrollment, error) {
	class, err := s.classRepo.GetByCode(ctx, classCode)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrClassNotFound
		}
		return nil, fmt.Errorf("failed to find class: %w", err)
	}

	// Check if already enrolled
	enrolled, err := s.enrollmentRepo.IsEnrolled(ctx, class.ID, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to check enrollment: %w", err)
	}
	if enrolled {
		return nil, ErrAlreadyEnrolled
	}

	// Create enrollment
	enrollment := &models.Enrollment{
		ID:         uuid.New(),
		ClassID:    class.ID,
		StudentID:  studentID,
		EnrolledAt: time.Now(),
	}

	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, ErrAlreadyEnrolled
		}
		return nil, fmt.Errorf("failed to create enrollment: %w", err)
	}

	return enrollment, nil
}

// GetStudentClasses returns all classes a student is enrolled in with class details.
func (s *enrollmentService) GetStudentClasses(
	ctx context.Context, studentID uuid.UUID,
) ([]models.EnrollmentWithClass, error) {
	classes, err := s.enrollmentRepo.GetClassesWithDetailsByStudentID(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student classes: %w", err)
	}

	return classes, nil
}

// GetClassStudents returns all students enrolled in a class with user details.
func (s *enrollmentService) GetClassStudents(
	ctx context.Context, classID uuid.UUID,
) ([]models.StudentInClass, error) {
	// Verify class exists
	_, err := s.classRepo.GetByID(ctx, classID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrClassNotFound
		}
		return nil, fmt.Errorf("failed to get class: %w", err)
	}

	students, err := s.enrollmentRepo.GetStudentsWithDetailsByClassID(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("failed to get class students: %w", err)
	}

	return students, nil
}

// Unenroll removes a student from a class.
func (s *enrollmentService) Unenroll(ctx context.Context, classID, studentID uuid.UUID) error {
	err := s.enrollmentRepo.Delete(ctx, classID, studentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotEnrolled
		}
		return fmt.Errorf("failed to unenroll student: %w", err)
	}

	return nil
}

// IsEnrolled checks if a student is enrolled in a class.
func (s *enrollmentService) IsEnrolled(ctx context.Context, classID, studentID uuid.UUID) (bool, error) {
	return s.enrollmentRepo.IsEnrolled(ctx, classID, studentID)
}
