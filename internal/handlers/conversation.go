package handlers

import (
	"strconv"

	"chat-assistant-backend/internal/errors"
	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/request"
	"chat-assistant-backend/internal/response"
	"chat-assistant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ConversationHandler handles conversation-related HTTP requests
type ConversationHandler struct {
	conversationService services.ConversationService
}

// NewConversationHandler creates a new conversation handler
func NewConversationHandler(conversationService services.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
	}
}

// GetConversations handles GET /api/v1/conversations
// @Summary Get Conversations
// @Description Retrieve conversations list with pagination
// @Tags Conversations
// @Accept json
// @Produce json
// @Param user_id query string true "User ID" Format(uuid)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse{data=response.ConversationListResponse} "Conversations list"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/conversations [get]
func (h *ConversationHandler) GetConversations(c *gin.Context) {
	// Parse user ID from query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		response.BadRequest(c, "MISSING_USER_ID", "User ID is required", "user_id query parameter is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_UUID", "Invalid user ID format", "User ID must be a valid UUID")
		return
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

	// Get conversations from service
	conversations, total, err := h.conversationService.GetConversationsByUserID(userID, page, limit)
	if err != nil {
		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to retrieve conversations")
		return
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Return success response
	conversationResponse := response.NewConversationListResponse(conversations)
	pagination := &response.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	response.SuccessPaginated(c, conversationResponse, pagination)
}

// GetConversation handles GET /api/v1/conversations/{id}
// @Summary Get Conversation
// @Description Retrieve a specific conversation by ID
// @Tags Conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID" Format(uuid)
// @Success 200 {object} response.Response{data=response.ConversationResponse} "Conversation details"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Conversation not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/conversations/{id} [get]
func (h *ConversationHandler) GetConversation(c *gin.Context) {
	// Parse conversation ID from path parameter
	conversationIDStr := c.Param("id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_UUID", "Invalid conversation ID format", "Conversation ID must be a valid UUID")
		return
	}

	// Get conversation from service
	conversation, err := h.conversationService.GetConversationByID(conversationID)
	if err != nil {
		if err == errors.ErrConversationNotFound {
			response.NotFound(c, "CONVERSATION_NOT_FOUND", "Conversation not found", "No conversation found with the specified ID")
			return
		}

		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to retrieve conversation")
		return
	}

	// Return success response
	conversationResponse := response.NewConversationResponse(conversation)
	response.Success(c, conversationResponse)
}

// DeleteConversation handles DELETE /api/v1/conversations/{id}
// @Summary Delete Conversation
// @Description Delete a specific conversation by ID
// @Tags Conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID" Format(uuid)
// @Success 200 {object} response.Response "Conversation deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Conversation not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/conversations/{id} [delete]
func (h *ConversationHandler) DeleteConversation(c *gin.Context) {
	// Parse conversation ID from path parameter
	conversationIDStr := c.Param("id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_UUID", "Invalid conversation ID format", "Conversation ID must be a valid UUID")
		return
	}

	// Delete conversation from service
	err = h.conversationService.DeleteConversation(conversationID)
	if err != nil {
		if err == errors.ErrConversationNotFound {
			response.NotFound(c, "CONVERSATION_NOT_FOUND", "Conversation not found", "No conversation found with the specified ID")
			return
		}

		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to delete conversation")
		return
	}

	// Return success response
	response.Success(c, gin.H{"message": "Conversation deleted successfully"})
}

// CreateConversation handles POST /api/v1/conversations
// @Summary Create Conversation
// @Description Create a new conversation with tags
// @Tags Conversations
// @Accept json
// @Produce json
// @Param conversation body request.CreateConversationRequest true "Conversation data"
// @Success 201 {object} response.Response{data=response.ConversationResponse} "Conversation created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/conversations [post]
func (h *ConversationHandler) CreateConversation(c *gin.Context) {
	var req request.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request data", err.Error())
		return
	}

	// 提取标签名称
	var tagNames []string
	for _, tag := range req.Tags {
		tagNames = append(tagNames, tag.Name)
	}

	// 创建对话模型
	conversation := &models.Conversation{
		UserID:      req.UserID,
		Title:       req.Title,
		Provider:    req.Provider,
		Model:       req.Model,
		SourceID:    req.SourceID,
		SourceTitle: req.SourceTitle,
	}

	// 创建对话和标签
	createdConversation, err := h.conversationService.CreateConversationWithTags(conversation, tagNames)
	if err != nil {
		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to create conversation")
		return
	}

	// Return success response
	conversationResponse := response.NewConversationResponse(createdConversation)
	response.Success(c, conversationResponse)
}

// UpdateConversationTags handles PUT /api/v1/conversations/{id}/tags
// @Summary Update Conversation Tags
// @Description Update tags for a specific conversation
// @Tags Conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID" Format(uuid)
// @Param tags body request.UpdateConversationTagsRequest true "Tags data"
// @Success 200 {object} response.Response "Tags updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Conversation not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/conversations/{id}/tags [put]
func (h *ConversationHandler) UpdateConversationTags(c *gin.Context) {
	// Parse conversation ID from path parameter
	conversationIDStr := c.Param("id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_UUID", "Invalid conversation ID format", "Conversation ID must be a valid UUID")
		return
	}

	var req request.UpdateConversationTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request data", err.Error())
		return
	}

	// 提取标签名称
	var tagNames []string
	for _, tag := range req.Tags {
		tagNames = append(tagNames, tag.Name)
	}

	// 更新对话标签
	err = h.conversationService.UpdateConversationTags(conversationID, tagNames)
	if err != nil {
		if err == errors.ErrConversationNotFound {
			response.NotFound(c, "CONVERSATION_NOT_FOUND", "Conversation not found", "No conversation found with the specified ID")
			return
		}

		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to update conversation tags")
		return
	}

	// Return success response
	response.Success(c, gin.H{"message": "Conversation tags updated successfully"})
}
