package handlers

import (
	"strconv"

	"chat-assistant-backend/internal/errors"
	"chat-assistant-backend/internal/response"
	"chat-assistant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MessageHandler handles message-related HTTP requests
type MessageHandler struct {
	messageService services.MessageService
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(messageService services.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// GetMessages handles GET /api/v1/messages
// @Summary Get Messages
// @Description Retrieve messages list with pagination
// @Tags Messages
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse{data=response.MessageListResponse} "Messages list"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/messages [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
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

	// Get messages from service
	messages, total, err := h.messageService.GetAllMessages(page, limit)
	if err != nil {
		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to retrieve messages")
		return
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Return success response
	messageResponse := response.NewMessageListResponse(messages)
	pagination := &response.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	response.SuccessPaginated(c, messageResponse, pagination)
}

// GetMessage handles GET /api/v1/messages/{id}
// @Summary Get Message
// @Description Retrieve a specific message by ID
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID" Format(uuid)
// @Success 200 {object} response.Response{data=response.MessageResponse} "Message details"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Message not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/messages/{id} [get]
func (h *MessageHandler) GetMessage(c *gin.Context) {
	// Parse message ID from path parameter
	messageIDStr := c.Param("id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_UUID", "Invalid message ID format", "Message ID must be a valid UUID")
		return
	}

	// Get message from service
	message, err := h.messageService.GetMessageByID(messageID)
	if err != nil {
		if err == errors.ErrMessageNotFound {
			response.NotFound(c, "MESSAGE_NOT_FOUND", "Message not found", "No message found with the specified ID")
			return
		}

		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to retrieve message")
		return
	}

	// Return success response
	messageResponse := response.NewMessageResponse(message)
	response.Success(c, messageResponse)
}

// DeleteMessage handles DELETE /api/v1/messages/{id}
// @Summary Delete Message
// @Description Delete a specific message by ID
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID" Format(uuid)
// @Success 200 {object} response.Response "Message deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Message not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	// Parse message ID from path parameter
	messageIDStr := c.Param("id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_UUID", "Invalid message ID format", "Message ID must be a valid UUID")
		return
	}

	// Delete message from service
	err = h.messageService.DeleteMessage(messageID)
	if err != nil {
		if err == errors.ErrMessageNotFound {
			response.NotFound(c, "MESSAGE_NOT_FOUND", "Message not found", "No message found with the specified ID")
			return
		}

		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to delete message")
		return
	}

	// Return success response
	response.Success(c, gin.H{"message": "Message deleted successfully"})
}

// GetConversationMessages handles GET /api/v1/conversations/{id}/messages
// @Summary Get Conversation Messages
// @Description Retrieve all messages in a specific conversation with pagination
// @Tags Conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID" Format(uuid)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse{data=response.MessageListResponse} "Messages list"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/conversations/{id}/messages [get]
func (h *MessageHandler) GetConversationMessages(c *gin.Context) {
	// Parse conversation ID from path parameter
	conversationIDStr := c.Param("id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_UUID", "Invalid conversation ID format", "Conversation ID must be a valid UUID")
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

	// Get messages from service
	messages, total, err := h.messageService.GetMessagesByConversationID(conversationID, page, limit)
	if err != nil {
		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to retrieve messages")
		return
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Return success response
	messageResponse := response.NewMessageListResponse(messages)
	pagination := &response.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	response.SuccessPaginated(c, messageResponse, pagination)
}
