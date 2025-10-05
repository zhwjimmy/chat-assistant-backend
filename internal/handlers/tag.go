package handlers

import (
	"chat-assistant-backend/internal/errors"
	"chat-assistant-backend/internal/request"
	"chat-assistant-backend/internal/response"
	"chat-assistant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TagHandler handles tag-related HTTP requests
type TagHandler struct {
	tagService services.TagService
}

// NewTagHandler creates a new tag handler
func NewTagHandler(tagService services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// GetTags handles GET /api/v1/tags
// @Summary Get Tags
// @Description Retrieve all tags
// @Tags Tags
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.TagListResponse} "Tags list"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/tags [get]
func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.tagService.GetAllTags()
	if err != nil {
		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to retrieve tags")
		return
	}

	tagResponse := response.NewTagListResponse(tags)
	response.Success(c, tagResponse)
}

// GetTag handles GET /api/v1/tags/{id}
// @Summary Get Tag
// @Description Retrieve a specific tag by ID
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID" Format(uuid)
// @Success 200 {object} response.Response{data=response.TagResponse} "Tag details"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Tag not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/tags/{id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	// Parse tag ID from path parameter
	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_UUID", "Invalid tag ID format", "Tag ID must be a valid UUID")
		return
	}

	// Get tag from service
	tag, err := h.tagService.GetTagByID(tagID)
	if err != nil {
		if err == errors.ErrTagNotFound {
			response.NotFound(c, "TAG_NOT_FOUND", "Tag not found", "No tag found with the specified ID")
			return
		}

		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to retrieve tag")
		return
	}

	// Return success response
	tagResponse := response.NewTagResponse(tag)
	response.Success(c, tagResponse)
}

// CreateTag handles POST /api/v1/tags
// @Summary Create Tag
// @Description Create a new tag
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag body request.CreateTagRequest true "Tag data"
// @Success 201 {object} response.Response{data=response.TagResponse} "Tag created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 409 {object} response.Response "Tag name already exists"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req request.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request data", err.Error())
		return
	}

	// Create tag
	tag, err := h.tagService.CreateTag(req.Name)
	if err != nil {
		if err == errors.ErrTagNameExists {
			response.Conflict(c, "TAG_NAME_EXISTS", "Tag name already exists", "A tag with this name already exists")
			return
		}

		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to create tag")
		return
	}

	// Return success response
	tagResponse := response.NewTagResponse(tag)
	response.Success(c, tagResponse)
}

// UpdateTag handles PUT /api/v1/tags/{id}
// @Summary Update Tag
// @Description Update an existing tag
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID" Format(uuid)
// @Param tag body request.UpdateTagRequest true "Tag data"
// @Success 200 {object} response.Response{data=response.TagResponse} "Tag updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Tag not found"
// @Failure 409 {object} response.Response "Tag name already exists"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	// Parse tag ID from path parameter
	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_UUID", "Invalid tag ID format", "Tag ID must be a valid UUID")
		return
	}

	var req request.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request data", err.Error())
		return
	}

	// Update tag
	tag, err := h.tagService.UpdateTag(tagID, req.Name)
	if err != nil {
		if err == errors.ErrTagNotFound {
			response.NotFound(c, "TAG_NOT_FOUND", "Tag not found", "No tag found with the specified ID")
			return
		}
		if err == errors.ErrTagNameExists {
			response.Conflict(c, "TAG_NAME_EXISTS", "Tag name already exists", "A tag with this name already exists")
			return
		}

		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to update tag")
		return
	}

	// Return success response
	tagResponse := response.NewTagResponse(tag)
	response.Success(c, tagResponse)
}

// DeleteTag handles DELETE /api/v1/tags/{id}
// @Summary Delete Tag
// @Description Delete a specific tag by ID
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID" Format(uuid)
// @Success 200 {object} response.Response "Tag deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Tag not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	// Parse tag ID from path parameter
	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_UUID", "Invalid tag ID format", "Tag ID must be a valid UUID")
		return
	}

	// Delete tag
	err = h.tagService.DeleteTag(tagID)
	if err != nil {
		if err == errors.ErrTagNotFound {
			response.NotFound(c, "TAG_NOT_FOUND", "Tag not found", "No tag found with the specified ID")
			return
		}

		response.InternalServerError(c, "INTERNAL_ERROR", "Internal server error", "Failed to delete tag")
		return
	}

	// Return success response
	response.Success(c, gin.H{"message": "Tag deleted successfully"})
}
