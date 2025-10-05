package request

// TagRequest represents a tag in API request
type TagRequest struct {
	ID   *string `json:"id,omitempty"` // 可选：现有标签ID
	Name string  `json:"name"`         // 必需：标签名称
}

// CreateTagRequest represents a request to create a tag
type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateTagRequest represents a request to update a tag
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateConversationTagsRequest represents a request to update conversation tags
type UpdateConversationTagsRequest struct {
	Tags []TagRequest `json:"tags" binding:"required"`
}
