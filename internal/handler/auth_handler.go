package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/tahiriqbal095/attendify/internal/models"
	"github.com/tahiriqbal095/attendify/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
	logger      zerolog.Logger
}

func NewAuthHandler(authService service.AuthService, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
		logger:      logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		BadRequest(c, "invalid request body")
		return
	}

	if err := h.validate.Struct(&input); err != nil {
		BadRequest(c, formatValidationError(err))
		return
	}

	if !input.Role.IsValid() {
		BadRequest(c, "role must be 'teacher' or 'student'")
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &input)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			Error(c, http.StatusConflict, "email already registered")
			return
		}
		h.logger.Error().Err(err).Str("email", input.Email).Msg("Failed to register user")
		InternalError(c)
		return
	}

	Success(c, http.StatusCreated, user.ToResponse())
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		BadRequest(c, "invalid request body")
		return
	}

	if err := h.validate.Struct(&input); err != nil {
		BadRequest(c, formatValidationError(err))
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			Unauthorized(c, "invalid email or password")
			return
		}
		h.logger.Error().Err(err).Str("email", input.Email).Msg("Failed to login user")
		InternalError(c)
		return
	}

	Success(c, http.StatusOK, response)
}

// formatValidationError converts validation errors to user-friendly messages.
func formatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				return e.Field() + " is required"
			case "email":
				return "invalid email format"
			case "min":
				return e.Field() + " is too short"
			case "max":
				return e.Field() + " is too long"
			case "oneof":
				return e.Field() + " must be one of: " + e.Param()
			default:
				return e.Field() + " is invalid"
			}
		}
	}
	return "validation failed"
}
