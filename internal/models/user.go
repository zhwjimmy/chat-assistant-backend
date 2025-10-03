package models

// User represents a user in the system
type User struct {
	Base
	Username string `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Avatar   string `json:"avatar" gorm:"size:255"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}
