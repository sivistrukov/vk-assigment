package validator

import (
	"reflect"
	"time"

	"github.com/go-playground/validator"
	"github.com/sivistrukov/vk-assigment/internal/models"
)

func New() *validator.Validate {
	validate := validator.New()
	_ = validate.RegisterValidation("dateValidation", dateValidation)
	_ = validate.RegisterValidation("sexValidation", sexValidation)

	return validate
}

func dateValidation(fl validator.FieldLevel) bool {
	value := fl.Field()
	var dateString string
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return false
		}
		dateString = value.Elem().String()
	} else {
		dateString = fl.Field().String()
	}

	layout := "02-01-2006"

	_, err := time.Parse(layout, dateString)
	return err == nil
}

func sexValidation(fl validator.FieldLevel) bool {
	value := fl.Field()
	var sexString string
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return false
		}
		sexString = value.Elem().String()
	} else {
		sexString = fl.Field().String()
	}

	return sexString == string(models.Male) || sexString == string(models.Female)
}
