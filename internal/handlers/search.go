package handlers

import (
	"strconv"
	"time"

	"chat-assistant-backend/internal/response"
	"chat-assistant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SearchHandler handles search-related HTTP requests
type SearchHandler struct {
	searchService services.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// Search handles GET /api/v1/search
// @Summary Search Conversations
// @Description Search conversations by title or message content, returns conversation list with matched messages
// @Tags Search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param user_id query string false "User ID" Format(uuid)
// @Param provider_id query string false "Provider ID (e.g., openai, gemini, claude)"
// @Param start_date query string false "Start date for filtering conversations" Format(date)
// @Param end_date query string false "End date for filtering conversations" Format(date)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse{data=response.SearchResponse} "Search results"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	// Parse search query
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "MISSING_QUERY", "Search query is required", "q query parameter is required")
		return
	}

	// Parse user ID (optional)
	var userID *uuid.UUID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = &parsed
		} else {
			response.BadRequest(c, "INVALID_UUID", "Invalid user ID format", "User ID must be a valid UUID")
			return
		}
	}

	// Parse provider ID (optional)
	var providerID *string
	if providerIDStr := c.Query("provider_id"); providerIDStr != "" {
		providerID = &providerIDStr
	}

	// Parse date range (optional)
	var startDate, endDate *time.Time
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		} else {
			response.BadRequest(c, "INVALID_DATE", "Invalid start date format", "Start date must be in YYYY-MM-DD format")
			return
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// 设置结束日期为当天的23:59:59
			endOfDay := parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			endDate = &endOfDay
		} else {
			response.BadRequest(c, "INVALID_DATE", "Invalid end date format", "End date must be in YYYY-MM-DD format")
			return
		}
	}

	// Parse pagination parameters
	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Perform search with matched messages
	searchResponse, total, err := h.searchService.SearchWithMatchedMessages(query, userID, providerID, startDate, endDate, page, limit)
	if err != nil {
		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to perform search")
		return
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Return success response
	pagination := &response.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	response.SuccessPaginated(c, searchResponse, pagination)
}
