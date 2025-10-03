package response

import (
	"time"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error information in response
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Response
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}

// MetaInfo represents metadata information
type MetaInfo struct {
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id,omitempty"`
	Version   string    `json:"version,omitempty"`
}

// MetaResponse represents a response with metadata
type MetaResponse struct {
	Response
	Meta *MetaInfo `json:"meta,omitempty"`
}
