package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"username,unique" binding:"required,uniqueField=username" json:"username"`
	Email string `gorm:"email,unique" binding:"required,uniqueField=email" json:"email"`
	Password string `gorm:"password" json:"password"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(err.Error())
	}

	u.Password = string(hash)

	return nil
}