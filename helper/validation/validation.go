package validation

import (
	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/yosikez/crudAuth/database"
	"github.com/yosikez/crudAuth/model"
	"gorm.io/gorm"
)

func RegisterCustomValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("uniqueField", UniqueField)
	}
}

func UniqueField(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	var user model.User
	field := fl.Param()

	result := database.DB.Table("users").Where(field+" = ?", value).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return true
	}

	parentIDField := fl.Parent().FieldByName("ID")
	if !parentIDField.IsValid() {
		return false
	}

	parentID := parentIDField.Interface().(uint)

	duplicateField := model.User{}
	duplicateResult := database.DB.Table("users").Where(field+" = ?", value).First(&duplicateField)
	if !errors.Is(duplicateResult.Error, gorm.ErrRecordNotFound) && duplicateField.Id != uint(parentID) {
		return false
	}

	return true
}
