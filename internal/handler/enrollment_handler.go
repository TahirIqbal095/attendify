package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/tahiriqbal095/attendify/internal/models"
	"github.com/tahiriqbal095/attendify/internal/service"
)

type EnrollmentHandler struct {
	enrollmentService service.EnrollmentService
	classService      service.ClassService
	logger            zerolog.Logger
	validate          *validator.Validate
}

func NewEnrollmentHandler(
	enrollmentService service.EnrollmentService,
	classService service.ClassService,
	logger zerolog.Logger,
) *EnrollmentHandler {
	return &EnrollmentHandler{
		enrollmentService: enrollmentService,
		classService:      classService,
		logger:            logger,
		validate:          validator.New(),
	}
}

// EnrollByCode handles POST /api/enrollments
// Students join a class using the class code.
func (h *EnrollmentHandler) EnrollByCode(c *gin.Context) {
	var input models.EnrollByCodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn().Err(err).Msg("invalid request body")
		BadRequest(c, "invalid request body")
		return
	}

	if err := h.validate.Struct(&input); err != nil {
		h.logger.Warn().Err(err).Msg("validation failed")
		BadRequest(c, formatValidationError(err))
		return
	}

	studentID, exists := c.Get("userID")
	if !exists {
		Unauthorized(c, "unauthorized")
		return
	}

	enrollment, err := h.enrollmentService.EnrollByCode(c.Request.Context(), input.ClassCode, studentID.(uuid.UUID))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrClassNotFound):
			NotFound(c, "class not found")
		case errors.Is(err, service.ErrAlreadyEnrolled):
			Error(c, http.StatusConflict, "already enrolled in this class")
		default:
			h.logger.Error().Err(err).Msg("failed to enroll student")
			InternalError(c)
		}
		return
	}

	h.logger.Info().
		Str("student_id", studentID.(uuid.UUID).String()).
		Str("class_id", enrollment.ClassID.String()).
		Msg("student enrolled successfully")

	Success(c, http.StatusCreated, enrollment)
}

// GetMyClasses handles GET /api/enrollments
// Returns all classes the authenticated student is enrolled in.
func (h *EnrollmentHandler) GetMyClasses(c *gin.Context) {
	studentID, exists := c.Get("userID")
	if !exists {
		Unauthorized(c, "unauthorized")
		return
	}

	classes, err := h.enrollmentService.GetStudentClasses(c.Request.Context(), studentID.(uuid.UUID))
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get enrolled classes")
		InternalError(c)
		return
	}

	Success(c, http.StatusOK, classes)
}

// GetClassStudents handles GET /api/classes/:id/students
// Returns all students enrolled in a class (teacher only).
func (h *EnrollmentHandler) GetClassStudents(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		BadRequest(c, "invalid class ID")
		return
	}

	// Verify teacher owns this class
	teacherID, exists := c.Get("userID")
	if !exists {
		Unauthorized(c, "unauthorized")
		return
	}

	class, err := h.classService.GetClass(c.Request.Context(), classID)
	if err != nil {
		if errors.Is(err, service.ErrClassNotFound) {
			NotFound(c, "class not found")
			return
		}
		h.logger.Error().Err(err).Msg("failed to get class")
		InternalError(c)
		return
	}

	if class.TeacherID != teacherID.(uuid.UUID) {
		Forbidden(c, "access denied")
		return
	}

	students, err := h.enrollmentService.GetClassStudents(c.Request.Context(), classID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get class students")
		InternalError(c)
		return
	}

	Success(c, http.StatusOK, students)
}

// Unenroll handles DELETE /api/enrollments/:classId
// Student leaves a class.
func (h *EnrollmentHandler) Unenroll(c *gin.Context) {
	classIDStr := c.Param("classId")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		BadRequest(c, "invalid class ID")
		return
	}

	studentID, exists := c.Get("userID")
	if !exists {
		Unauthorized(c, "unauthorized")
		return
	}

	err = h.enrollmentService.Unenroll(c.Request.Context(), classID, studentID.(uuid.UUID))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotEnrolled):
			NotFound(c, "not enrolled in this class")
		default:
			h.logger.Error().Err(err).Msg("failed to unenroll student")
			InternalError(c)
		}
		return
	}

	h.logger.Info().
		Str("student_id", studentID.(uuid.UUID).String()).
		Str("class_id", classID.String()).
		Msg("student unenrolled successfully")

	Success(c, http.StatusOK, nil)
}
