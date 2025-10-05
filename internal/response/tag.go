package response

import (
	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
)

type TagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// NewTagResponse creates a TagResponse from models.Tag
func NewTagResponse(tag *models.Tag) *TagResponse {
	return &TagResponse{
		ID:   tag.Base.ID,
		Name: tag.Name,
	}
}

// TagListResponse represents a list of tags in API response
type TagListResponse struct {
	Tags []TagResponse `json:"tags"`
}

// NewTagListResponse creates a TagListResponse from a slice of models.Tag
func NewTagListResponse(tags []*models.Tag) *TagListResponse {
	tagResponses := make([]TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = *NewTagResponse(tag)
	}

	return &TagListResponse{
		Tags: tagResponses,
	}
}
