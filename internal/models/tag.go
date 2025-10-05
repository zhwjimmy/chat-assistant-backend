package models

type Tag struct {
	Base
	Name string `gorm:"type:varchar(500);not null" json:"name"`
}

// TableName returns the table name for the Conversation model
func (Tag) TableName() string {
	return "tags"
}
