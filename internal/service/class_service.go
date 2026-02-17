package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tahiriqbal095/attendify/internal/models"
	"github.com/tahiriqbal095/attendify/internal/repository"
)

type ClassService interface {
	CreateClass(ctx context.Context, teacherID uuid.UUID, input *models.CreateClassInput) (*models.Class, error)
	GetClass(ctx context.Context, id uuid.UUID) (*models.Class, error)
	GetClassByCode(ctx context.Context, code string) (*models.Class, error)
	GetTeacherClasses(ctx context.Context, teacherID uuid.UUID) ([]models.Class, error)
	DeleteClass(ctx context.Context, teacherID, classID uuid.UUID) error
}

type classService struct {
	classRepo repository.ClassRepository
}

func NewClassService(classRepo repository.ClassRepository) ClassService {
	return &classService{classRepo: classRepo}
}

func (s *classService) CreateClass(ctx context.Context, teacherID uuid.UUID, input *models.CreateClassInput) (*models.Class, error) {
	code, err := s.generateUniqueCode(ctx)
	if err != nil {
		return nil, err
	}

	class := &models.Class{
		ID:        uuid.New(),
		Name:      input.Name,
		Code:      code,
		TeacherID: teacherID,
		CreatedAt: time.Now(),
	}

	if err := s.classRepo.Create(ctx, class); err != nil {
		return nil, fmt.Errorf("failed to create class: %w", err)
	}

	return class, nil
}

func (s *classService) GetClass(ctx context.Context, id uuid.UUID) (*models.Class, error) {
	class, err := s.classRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrClassNotFound
		}
		return nil, fmt.Errorf("failed to get class: %w", err)
	}

	return class, nil
}

func (s *classService) GetClassByCode(ctx context.Context, code string) (*models.Class, error) {
	class, err := s.classRepo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrClassNotFound
		}
		return nil, fmt.Errorf("failed to get class by code: %w", err)
	}

	return class, nil
}

func (s *classService) GetTeacherClasses(ctx context.Context, teacherID uuid.UUID) ([]models.Class, error) {
	classes, err := s.classRepo.GetByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher classes: %w", err)
	}

	return classes, nil
}

func (s *classService) DeleteClass(ctx context.Context, teacherID, classID uuid.UUID) error {
	class, err := s.classRepo.GetByID(ctx, classID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrClassNotFound
		}
		return fmt.Errorf("failed to get class: %w", err)
	}

	if class.TeacherID != teacherID {
		return ErrNotClassOwner
	}

	if err := s.classRepo.Delete(ctx, classID); err != nil {
		return fmt.Errorf("failed to delete class: %w", err)
	}

	return nil
}

func (s *classService) generateUniqueCode(ctx context.Context) (string, error) {
	const maxAttempts = 5

	for i := 0; i < maxAttempts; i++ {
		code, err := generateCode(6)
		if err != nil {
			continue
		}

		_, err = s.classRepo.GetByCode(ctx, code)
		if errors.Is(err, repository.ErrNotFound) {
			return code, nil
		}
	}

	return "", ErrCodeGeneration
}

func generateCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	code := base32.StdEncoding.EncodeToString(bytes)
	return strings.ToUpper(code[:length]), nil
}
