package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// TODO: implement  better way of handling the replies of comment
type Comment struct {
	ID              string    `gorm:"type:text;primaryKey" json:"id" binding:"-"`
	Content         string    `gorm:"type:text;not null" binding:"required" validate:"required" json:"content"`
	ActivityID      uint      `gorm:"type:text;foreignKey:ActivityID;not null" binding:"-" validate:"required" json:"activityID"`
	ParentID        *string   `gorm:"type:text" json:"parentId"`
	Replies         []Comment `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"replies" binding:"-"`
	CreatedByHandle string    `gorm:"not null;type:text" validate:"required" json:"-" binding:"-"`
	CreatedBy       User      `gorm:"references:Handle;foreignKey:CreatedByHandle" json:"createdBy" binding:"-"`
	CreatedAt       time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt       time.Time `json:"-"`
}

// Constructor and Initialization Methods

// NewComment creates a new Comment instance
func NewComment(parentID *string, activityID uint, content string, createdBy User) *Comment {
	return &Comment{
		ID:              generateID(),
		ParentID:        parentID,
		ActivityID:      activityID,
		CreatedByHandle: createdBy.Handle,
		CreatedBy:       createdBy,
		Content:         content,
	}
}

// FromBinding initializes a Comment instance from binding data
func (c *Comment) FromBinding(CreatedBy User, activityID uint) {

	c.ID = generateID()
	c.CreatedByHandle = CreatedBy.Handle
	c.CreatedBy = CreatedBy
	c.ActivityID = activityID
}

// Helper function -------------------------------------------------

// generateId generates a new id for Comment
func generateID() string {
	return "CMT" + strings.ToUpper(uuid.NewString()[:5])
}
