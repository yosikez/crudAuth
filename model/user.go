package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id       uint      `gorm:"column:id" json:"id"`
	Username string    `gorm:"column:username;unique" binding:"required,uniqueField=username" json:"username"`
	Email    string    `gorm:"column:email;unique" binding:"required,uniqueField=email" json:"email"`
	Password string    `gorm:"column:password" json:"password"`
	CreateAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdateAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(err.Error())
	}

	u.Password = string(hash)
	u.CreateAt = now
	u.UpdateAt = now

	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdateAt = time.Now()
	return nil
}
