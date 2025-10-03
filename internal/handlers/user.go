package handlers

import (
	"chat-assistant-backend/internal/errors"
	"chat-assistant-backend/internal/response"
	"chat-assistant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser handles GET /api/v1/users/{id}
// @Summary Get User
// @Description Retrieve a specific user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID" Format(uuid)
// @Success 200 {object} response.Response{data=response.UserResponse} "User details"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "User not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	// Parse user ID from path parameter
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_UUID", "Invalid user ID format", "User ID must be a valid UUID")
		return
	}

	// Get user from service
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if err == errors.ErrUserNotFound {
			response.NotFound(c, "USER_NOT_FOUND", "User not found", "No user found with the specified ID")
			return
		}

		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to retrieve user")
		return
	}

	// Return success response
	userResponse := response.NewUserResponse(user)
	response.Success(c, userResponse)
}
