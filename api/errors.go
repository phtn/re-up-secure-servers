package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents the structure of error responses sent to clients
type ErrorResponse struct {
	Status           int               `json:"status"`                      // HTTP status code
	Code             string            `json:"code"`                        // Internal error code
	Message          string            `json:"message"`                     // User-friendly error message
	Details          interface{}       `json:"details,omitempty"`           // Additional error details
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"` // Validation errors
}

// ValidationError represents individual field validation errors
type ValidationError struct {
	Field   string `json:"field"`   // Field that failed validation
	Message string `json:"message"` // Error message for this field
}

// AppError is our custom error type for application errors
type AppError struct {
	Status  int
	Code    string
	Message string
	Err     error
	Details interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Predefined error codes
const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeLocked       = "ACCOUNT_LOCKED"
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrCodeBadRequest   = "BAD_REQUEST"
	ErrCodeDuplicate    = "DUPLICATE_ENTRY"
	ErrCodeRateLimit    = "RATE_LIMIT_EXCEEDED"
	ErrCodeNotFound     = "NOT_FOUND"
)

// Common application errors
var (
	ErrInvalidInput = &AppError{
		Status:  fiber.StatusBadRequest,
		Code:    ErrCodeBadRequest,
		Message: "Invalid input provided",
	}

	ErrLocked = &AppError{
		Status:  fiber.StatusLocked,
		Code:    ErrCodeLocked,
		Message: "Mother account inactive",
	}

	ErrUnauthorized = &AppError{
		Status:  fiber.StatusUnauthorized,
		Code:    ErrCodeUnauthorized,
		Message: "Unauthorized access",
	}

	ErrBadRequest = &AppError{
		Status:  fiber.StatusBadRequest,
		Code:    ErrCodeBadRequest,
		Message: "Invalid Body Params",
	}

	ErrInternalServer = &AppError{
		Status:  fiber.StatusInternalServerError,
		Code:    ErrCodeInternal,
		Message: "Internal server error",
	}
	ErrNotFound = &AppError{
		Status:  fiber.StatusNotFound,
		Code:    ErrCodeNotFound,
		Message: "Not Found",
	}
)

// ErrorHandler is a middleware that catches all errors and formats them appropriately
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to internal server error
	status := fiber.StatusInternalServerError
	code := ErrCodeInternal
	message := "An unexpected error occurred"
	var details interface{}
	var validationErrors []ValidationError

	// Handle different error types
	switch e := err.(type) {
	case *AppError:
		status = e.Status
		code = e.Code
		message = e.Message
		details = e.Details

	case *fiber.Error:
		status = e.Code
		message = e.Message

	case ValidationErrors:
		status = fiber.StatusBadRequest
		code = ErrCodeValidation
		message = "Validation failed"
		validationErrors = e

	}

	// Don't expose internal error details in production
	if fiber.IsChild() { // IsChild() returns true in production
		if status == fiber.StatusInternalServerError {
			message = "An unexpected error occurred"
			details = nil
		}
	}
	return c.Status(status).JSON(ErrorResponse{
		Status:           status,
		Code:             code,
		Message:          message,
		Details:          details,
		ValidationErrors: validationErrors,
	})
}

func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// NewAppError creates a new application error
func NewAppError(status int, code, message string, err error) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

type ValidationErrors []ValidationError

// Error implements the error interface for ValidationErrors
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}
	return fmt.Sprintf("validation failed: %s", ve[0].Message)
}

// Helper function to create validation errors
func NewValidationErrors(errors ...ValidationError) ValidationErrors {
	return ValidationErrors(errors)
}

// WithDetails adds details to an AppError
func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
}
