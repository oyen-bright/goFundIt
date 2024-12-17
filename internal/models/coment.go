package models

type Comment struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Content    string `gorm:"type:text;not null" binding:"required" validate:"required" json:"content"`
	ActivityID string `gorm:"type:text;foreignKey:ActivityID;not null" binding:"required" validate:"required" json:"activityID"`
	CreatedAt  string `gorm:"not null" json:"createdAt"`
	UpdatedAt  string `json:"-"`
}
