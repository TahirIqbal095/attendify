package models

import (
	"time"

	"github.com/google/uuid"
)

type Enrollment struct {
	ID         uuid.UUID `json:"id"`
	ClassID    uuid.UUID `json:"class_id"`
	StudentID  uuid.UUID `json:"student_id"`
	EnrolledAt time.Time `json:"enrolled_at"`
}

type EnrollByCodeInput struct {
	ClassCode string `json:"class_code" validate:"required,min=4,max=10"`
}

type EnrollmentResponse struct {
	ID         uuid.UUID `json:"id"`
	ClassID    uuid.UUID `json:"class_id"`
	StudentID  uuid.UUID `json:"student_id"`
	EnrolledAt time.Time `json:"enrolled_at"`
}

type EnrollmentWithClass struct {
	ID         uuid.UUID     `json:"id"`
	Class      ClassResponse `json:"class"`
	EnrolledAt time.Time     `json:"enrolled_at"`
}

type StudentInClass struct {
	ID         uuid.UUID    `json:"id"`
	Student    UserResponse `json:"student"`
	EnrolledAt time.Time    `json:"enrolled_at"`
}

func (e *Enrollment) ToResponse() EnrollmentResponse {
	return EnrollmentResponse{
		ID:         e.ID,
		ClassID:    e.ClassID,
		StudentID:  e.StudentID,
		EnrolledAt: e.EnrolledAt,
	}
}
