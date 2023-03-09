package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	Id          uint      `gorm:"column:id" json:"id"`
	Title       string    `gorm:"column:title" json:"title" binding:"required"`
	Description string    `gorm:"column:description;type:text" json:"description" binding:"required"`
	DueDate     string    `gorm:"column:due_date" json:"due_date" binding:"required,datetime=2006-01-02"`
	IsComplete  bool      `gorm:"column:is_complete;default:false" json:"is_complete"`
	UserId      uint      `gorm:"foreignKey:User;OnUpdate:CASCADE;OnDelete:CASCADE" json:"user_id"`
	CreateAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdateAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (t *Todo) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	t.CreateAt = now
	t.UpdateAt = now
	return nil
}

func (t *Todo) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateAt = time.Now()
	return nil
}
