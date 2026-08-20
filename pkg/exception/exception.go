package exception

import (
	"fmt"
	"time"
)

type AppError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Detail    string    `json:"detail"`
	Timestamp time.Time `json:"timestamp"`
	Field     string    `json:"field,omitempty"`
}

func (e *AppError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("exception %d: %s - field: %s (%s)", e.Code, e.Message, e.Field, e.Detail)
	}
	return fmt.Sprintf("exception %d: %s (%s)", e.Code, e.Message, e.Detail)
}

// NewBadRequest HTTP Status Code Errors
func NewBadRequest(message, detail string) *AppError {
	return &AppError{
		Code:      400,
		Message:   message,
		Detail:    detail,
		Timestamp: time.Now(),
	}
}

func NewUnauthorized(message string) *AppError {
	return &AppError{
		Code:      401,
		Message:   message,
		Detail:    "invalid or missing token",
		Timestamp: time.Now(),
	}
}

func NewForbidden(message string) *AppError {
	return &AppError{
		Code:      403,
		Message:   message,
		Detail:    "insufficient permissions",
		Timestamp: time.Now(),
	}
}

func NewNotFound(message, detail string) *AppError {
	return &AppError{
		Code:      404,
		Message:   message,
		Detail:    detail,
		Timestamp: time.Now(),
	}
}

func NewConflict(message, detail string) *AppError {
	return &AppError{
		Code:      409,
		Message:   message,
		Detail:    detail,
		Timestamp: time.Now(),
	}
}

func NewTooManyRequest(message string) *AppError {
	return &AppError{
		Code:      429,
		Message:   message,
		Detail:    "rate limit exceeded",
		Timestamp: time.Now(),
	}
}

func NewInternal(message string) *AppError {
	return &AppError{
		Code:      500,
		Message:   message,
		Detail:    "internal server exception",
		Timestamp: time.Now(),
	}
}

func NewServiceUnavailable(message string) *AppError {
	return &AppError{
		Code:      503,
		Message:   message,
		Detail:    "service temporarily unavailable",
		Timestamp: time.Now(),
	}
}

// Validation Errors
func NewValidationError(field, message string) *AppError {
	return &AppError{
		Code:      422,
		Message:   message,
		Detail:    fmt.Sprintf("validation failed for field: %s", field),
		Field:     field,
		Timestamp: time.Now(),
	}
}

// Business Logic Errors
func NewBusinessError(message, detail string) *AppError {
	return &AppError{
		Code:      422,
		Message:   message,
		Detail:    detail,
		Timestamp: time.Now(),
	}
}

// Database Errors
func NewDatabaseError(message string) *AppError {
	return &AppError{
		Code:      500,
		Message:   message,
		Detail:    "database operation failed",
		Timestamp: time.Now(),
	}
}

// External Service Errors
func NewExternalServiceError(service, message string) *AppError {
	return &AppError{
		Code:      502,
		Message:   message,
		Detail:    fmt.Sprintf("external service error: %s", service),
		Timestamp: time.Now(),
	}
}

// Custom error with specific code
func NewCustomError(code int, message, detail string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Detail:    detail,
		Timestamp: time.Now(),
	}
}

// Helper function to check if error is AppError
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}
