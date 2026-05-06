package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(data interface{}) error {
	return validate.Struct(data)
}

func ValidationErrorResponse(c *gin.Context, err error) {
	var errors map[string]string

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		errors = make(map[string]string)
		for _, ve := range validationErrs {
			field := strings.ToLower(ve.Field())
			switch ve.Tag() {
			case "required":
				errors[field] = "Field ini wajib diisi"
			case "email":
				errors[field] = "Format email tidak valid"
			case "min":
				errors[field] = "Minimal " + ve.Param() + " karakter"
			case "max":
				errors[field] = "Maksimal " + ve.Param() + " karakter"
			case "numeric":
				errors[field] = "Harus berupa angka"
			default:
				errors[field] = "Field tidak valid"
			}
		}
	} else {
		errors = map[string]string{"error": err.Error()}
	}

	ErrorResponse(c, 400, "Validation failed", nil)
	c.JSON(400, gin.H{
		"success": false,
		"message": "Validation failed",
		"errors":  errors,
	})
}