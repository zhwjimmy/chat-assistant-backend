package response

import (
	"net/http"

	"chat-assistant-backend/internal/errors"

	"github.com/gin-gonic/gin"
)

// Success sends a success response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMeta sends a success response with metadata
func SuccessWithMeta(c *gin.Context, data interface{}, meta *MetaInfo) {
	c.JSON(http.StatusOK, MetaResponse{
		Response: Response{
			Success: true,
			Data:    data,
		},
		Meta: meta,
	})
}

// SuccessPaginated sends a paginated success response
func SuccessPaginated(c *gin.Context, data interface{}, pagination *PaginationInfo) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Response: Response{
			Success: true,
			Data:    data,
		},
		Pagination: pagination,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code, message, details string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// AppError sends an error response from AppError
func AppError(c *gin.Context, err *errors.AppError) {
	c.JSON(err.Status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    err.Code,
			Message: err.Message,
			Details: err.Details,
		},
	})
}

// BadRequest sends a bad request response
func BadRequest(c *gin.Context, code, message, details string) {
	Error(c, http.StatusBadRequest, code, message, details)
}

// Unauthorized sends an unauthorized response
func Unauthorized(c *gin.Context, code, message, details string) {
	Error(c, http.StatusUnauthorized, code, message, details)
}

// Forbidden sends a forbidden response
func Forbidden(c *gin.Context, code, message, details string) {
	Error(c, http.StatusForbidden, code, message, details)
}

// NotFound sends a not found response
func NotFound(c *gin.Context, code, message, details string) {
	Error(c, http.StatusNotFound, code, message, details)
}

// Conflict sends a conflict response
func Conflict(c *gin.Context, code, message, details string) {
	Error(c, http.StatusConflict, code, message, details)
}

// InternalServerError sends an internal server error response
func InternalServerError(c *gin.Context, code, message, details string) {
	Error(c, http.StatusInternalServerError, code, message, details)
}

// ServiceUnavailable sends a service unavailable response
func ServiceUnavailable(c *gin.Context, code, message, details string) {
	Error(c, http.StatusServiceUnavailable, code, message, details)
}
