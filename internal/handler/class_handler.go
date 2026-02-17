package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/tahiriqbal095/attendify/internal/middleware"
	"github.com/tahiriqbal095/attendify/internal/models"
	"github.com/tahiriqbal095/attendify/internal/service"
)

type ClassHandler struct {
	classService service.ClassService
	validate     *validator.Validate
	logger       zerolog.Logger
}

func NewClassHandler(classService service.ClassService, logger zerolog.Logger) *ClassHandler {
	return &ClassHandler{
		classService: classService,
		validate:     validator.New(),
		logger:       logger,
	}
}

func (h *ClassHandler) Create(c *gin.Context) {
	var input models.CreateClassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		BadRequest(c, "invalid request body")
		return
	}

	if err := h.validate.Struct(&input); err != nil {
		BadRequest(c, formatValidationError(err))
		return
	}

	teacherID := middleware.GetUserID(c)
	class, err := h.classService.CreateClass(c.Request.Context(), teacherID, &input)
	if err != nil {
		h.logger.Error().Err(err).Str("teacher_id", teacherID.String()).Msg("Failed to create class")
		InternalError(c)
		return
	}

	Success(c, http.StatusCreated, class.ToResponse())
}

func (h *ClassHandler) List(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	classes, err := h.classService.GetTeacherClasses(c.Request.Context(), teacherID)
	if err != nil {
		h.logger.Error().Err(err).Str("teacher_id", teacherID.String()).Msg("Failed to list classes")
		InternalError(c)
		return
	}

	// Convert to response format
	response := make([]models.ClassResponse, len(classes))
	for i, class := range classes {
		response[i] = class.ToResponse()
	}

	Success(c, http.StatusOK, response)
}

func (h *ClassHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	classID, err := uuid.Parse(idParam)
	if err != nil {
		BadRequest(c, "invalid class id")
		return
	}

	class, err := h.classService.GetClass(c.Request.Context(), classID)
	if err != nil {
		if errors.Is(err, service.ErrClassNotFound) {
			NotFound(c, "class not found")
			return
		}
		h.logger.Error().Err(err).Str("class_id", idParam).Msg("Failed to get class")
		InternalError(c)
		return
	}

	Success(c, http.StatusOK, class.ToResponse())
}

func (h *ClassHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	classID, err := uuid.Parse(idParam)
	if err != nil {
		BadRequest(c, "invalid class id")
		return
	}

	teacherID := middleware.GetUserID(c)
	err = h.classService.DeleteClass(c.Request.Context(), teacherID, classID)
	if err != nil {
		if errors.Is(err, service.ErrClassNotFound) {
			NotFound(c, "class not found")
			return
		}
		if errors.Is(err, service.ErrNotClassOwner) {
			Forbidden(c, "not the owner of this class")
			return
		}
		h.logger.Error().Err(err).Str("class_id", idParam).Msg("Failed to delete class")
		InternalError(c)
		return
	}

	Success(c, http.StatusOK, gin.H{"message": "class deleted"})
}
