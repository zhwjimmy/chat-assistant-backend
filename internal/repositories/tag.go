package repositories

import (
	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TagRepository defines the interface for tag repository
type TagRepository interface {
	GetByID(id uuid.UUID) (*models.Tag, error)
	GetByName(name string) (*models.Tag, error)
	GetByNames(names []string) ([]*models.Tag, error)
	Create(tag *models.Tag) error
	Update(tag *models.Tag) error
	Delete(id uuid.UUID) error
	FindAll() ([]*models.Tag, error)
	CreateOrGetTags(names []string) ([]*models.Tag, error)
}

// TagRepositoryImpl handles tag data access
type TagRepositoryImpl struct {
	db *gorm.DB
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *gorm.DB) TagRepository {
	return &TagRepositoryImpl{
		db: db,
	}
}

// GetByID retrieves a tag by ID
func (r *TagRepositoryImpl) GetByID(id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("id = ?", id).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil tag and nil error for not found
		}
		return nil, err
	}
	return &tag, nil
}

// GetByName retrieves a tag by name
func (r *TagRepositoryImpl) GetByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil tag and nil error for not found
		}
		return nil, err
	}
	return &tag, nil
}

// GetByNames retrieves tags by their names
func (r *TagRepositoryImpl) GetByNames(names []string) ([]*models.Tag, error) {
	if len(names) == 0 {
		return []*models.Tag{}, nil
	}

	var tags []*models.Tag
	err := r.db.Where("name IN ?", names).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// Create creates a new tag
func (r *TagRepositoryImpl) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

// Update updates an existing tag
func (r *TagRepositoryImpl) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

// Delete soft deletes a tag by ID
func (r *TagRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Tag{}, id).Error
}

// FindAll retrieves all tags
func (r *TagRepositoryImpl) FindAll() ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.Order("name ASC").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// CreateOrGetTags creates new tags or returns existing ones by names
func (r *TagRepositoryImpl) CreateOrGetTags(names []string) ([]*models.Tag, error) {
	if len(names) == 0 {
		return []*models.Tag{}, nil
	}

	// 去重
	uniqueNames := make(map[string]bool)
	var uniqueNameList []string
	for _, name := range names {
		if !uniqueNames[name] {
			uniqueNames[name] = true
			uniqueNameList = append(uniqueNameList, name)
		}
	}

	// 获取已存在的标签
	existingTags, err := r.GetByNames(uniqueNameList)
	if err != nil {
		return nil, err
	}

	// 创建已存在标签的映射
	existingTagMap := make(map[string]*models.Tag)
	for _, tag := range existingTags {
		existingTagMap[tag.Name] = tag
	}

	// 找出需要创建的标签
	var tagsToCreate []*models.Tag
	for _, name := range uniqueNameList {
		if _, exists := existingTagMap[name]; !exists {
			tagsToCreate = append(tagsToCreate, &models.Tag{
				Name: name,
			})
		}
	}

	// 批量创建新标签
	if len(tagsToCreate) > 0 {
		err = r.db.Create(&tagsToCreate).Error
		if err != nil {
			return nil, err
		}
	}

	// 合并结果
	var result []*models.Tag
	result = append(result, existingTags...)
	result = append(result, tagsToCreate...)

	return result, nil
}
