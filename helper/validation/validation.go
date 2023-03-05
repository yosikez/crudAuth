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

	duplicateEmployee := model.User{}
	duplicateResult := database.DB.Table("users").Where(field+" = ?", value).First(&duplicateEmployee)
	if !errors.Is(duplicateResult.Error, gorm.ErrRecordNotFound) && duplicateEmployee.ID != uint(parentID) {
		return false
	}

	return true
}