package models

type Tag struct {
	Base
	Name string `gorm:"type:varchar(500);not null" json:"name"`
}

// TableName returns the table name for the Tag model
func (Tag) TableName() string {
	return "tags"
}

// ToESDocument converts Tag to TagDocument for Elasticsearch
func (t *Tag) ToESDocument() TagDocument {
	return TagDocument{
		ID:        t.Base.ID,
		Name:      t.Name,
		CreatedAt: t.Base.CreatedAt,
		UpdatedAt: t.Base.UpdatedAt,
	}
}
