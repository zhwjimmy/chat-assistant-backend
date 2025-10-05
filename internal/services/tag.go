package services

import (
	"chat-assistant-backend/internal/errors"
	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/repositories"

	"github.com/google/uuid"
)

// TagService defines the interface for tag service
type TagService interface {
	GetTagByID(id uuid.UUID) (*models.Tag, error)
	GetTagByName(name string) (*models.Tag, error)
	GetAllTags() ([]*models.Tag, error)
	CreateTag(name string) (*models.Tag, error)
	UpdateTag(id uuid.UUID, name string) (*models.Tag, error)
	DeleteTag(id uuid.UUID) error
	CreateOrGetTags(names []string) ([]*models.Tag, error)
}

// TagServiceImpl handles tag business logic
type TagServiceImpl struct {
	tagRepo repositories.TagRepository
}

// NewTagService creates a new tag service
func NewTagService(tagRepo repositories.TagRepository) TagService {
	return &TagServiceImpl{
		tagRepo: tagRepo,
	}
}

// GetTagByID retrieves a tag by ID
func (s *TagServiceImpl) GetTagByID(id uuid.UUID) (*models.Tag, error) {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if tag == nil {
		return nil, errors.ErrTagNotFound
	}

	return tag, nil
}

// GetTagByName retrieves a tag by name
func (s *TagServiceImpl) GetTagByName(name string) (*models.Tag, error) {
	tag, err := s.tagRepo.GetByName(name)
	if err != nil {
		return nil, err
	}

	if tag == nil {
		return nil, errors.ErrTagNotFound
	}

	return tag, nil
}

// GetAllTags retrieves all tags
func (s *TagServiceImpl) GetAllTags() ([]*models.Tag, error) {
	return s.tagRepo.FindAll()
}

// CreateTag creates a new tag
func (s *TagServiceImpl) CreateTag(name string) (*models.Tag, error) {
	// 检查标签是否已存在
	existingTag, err := s.tagRepo.GetByName(name)
	if err != nil {
		return nil, err
	}

	if existingTag != nil {
		return existingTag, nil // 返回已存在的标签
	}

	// 创建新标签
	tag := &models.Tag{
		Name: name,
	}

	err = s.tagRepo.Create(tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// UpdateTag updates an existing tag
func (s *TagServiceImpl) UpdateTag(id uuid.UUID, name string) (*models.Tag, error) {
	// 检查标签是否存在
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if tag == nil {
		return nil, errors.ErrTagNotFound
	}

	// 检查新名称是否已被其他标签使用
	existingTag, err := s.tagRepo.GetByName(name)
	if err != nil {
		return nil, err
	}

	if existingTag != nil && existingTag.ID != id {
		return nil, errors.ErrTagNameExists
	}

	// 更新标签
	tag.Name = name
	err = s.tagRepo.Update(tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// DeleteTag deletes a tag by ID
func (s *TagServiceImpl) DeleteTag(id uuid.UUID) error {
	// 检查标签是否存在
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return err
	}

	if tag == nil {
		return errors.ErrTagNotFound
	}

	// 删除标签
	return s.tagRepo.Delete(id)
}

// CreateOrGetTags creates new tags or returns existing ones by names
func (s *TagServiceImpl) CreateOrGetTags(names []string) ([]*models.Tag, error) {
	if len(names) == 0 {
		return []*models.Tag{}, nil
	}

	return s.tagRepo.CreateOrGetTags(names)
}
