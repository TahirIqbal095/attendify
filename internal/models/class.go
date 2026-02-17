package models

import (
	"time"

	"github.com/google/uuid"
)

type Class struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	TeacherID uuid.UUID `json:"teacher_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateClassInput struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type ClassResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	TeacherID uuid.UUID `json:"teacher_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Class) ToResponse() ClassResponse {
	return ClassResponse{
		ID:        c.ID,
		Name:      c.Name,
		Code:      c.Code,
		TeacherID: c.TeacherID,
		CreatedAt: c.CreatedAt,
	}
}
