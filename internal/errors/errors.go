package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError creates a new application error
func NewAppError(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// Predefined error codes
const (
	// General errors
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeBadRequest   = "BAD_REQUEST"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
	ErrCodeConflict     = "CONFLICT"
	ErrCodeValidation   = "VALIDATION_ERROR"

	// Database errors
	ErrCodeDBConnection = "DB_CONNECTION_ERROR"
	ErrCodeDBQuery      = "DB_QUERY_ERROR"
	ErrCodeDBMigration  = "DB_MIGRATION_ERROR"

	// Configuration errors
	ErrCodeConfigLoad = "CONFIG_LOAD_ERROR"
)

// Predefined errors
var (
	ErrInternal     = NewAppError(ErrCodeInternal, "Internal server error", http.StatusInternalServerError)
	ErrNotFound     = NewAppError(ErrCodeNotFound, "Resource not found", http.StatusNotFound)
	ErrBadRequest   = NewAppError(ErrCodeBadRequest, "Bad request", http.StatusBadRequest)
	ErrUnauthorized = NewAppError(ErrCodeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	ErrForbidden    = NewAppError(ErrCodeForbidden, "Forbidden", http.StatusForbidden)
	ErrConflict     = NewAppError(ErrCodeConflict, "Resource conflict", http.StatusConflict)
	ErrValidation   = NewAppError(ErrCodeValidation, "Validation error", http.StatusBadRequest)

	ErrDBConnection = NewAppError(ErrCodeDBConnection, "Database connection error", http.StatusInternalServerError)
	ErrDBQuery      = NewAppError(ErrCodeDBQuery, "Database query error", http.StatusInternalServerError)
	ErrDBMigration  = NewAppError(ErrCodeDBMigration, "Database migration error", http.StatusInternalServerError)

	ErrConfigLoad = NewAppError(ErrCodeConfigLoad, "Configuration load error", http.StatusInternalServerError)
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *AppError   `json:"error,omitempty"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(err *AppError) *Response {
	return &Response{
		Success: false,
		Error:   err,
	}
}
