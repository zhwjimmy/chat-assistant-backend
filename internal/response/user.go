package response

import (
	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
)

// UserResponse represents a user in API response
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// UserListResponse represents a list of users in API response
type UserListResponse struct {
	Users []UserResponse `json:"users"`
}

// NewUserResponse creates a UserResponse from models.User
func NewUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.Base.ID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		CreatedAt: user.Base.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.Base.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// NewUserListResponse creates a UserListResponse from a slice of models.User
func NewUserListResponse(users []*models.User) *UserListResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *NewUserResponse(user)
	}

	return &UserListResponse{
		Users: userResponses,
	}
}
