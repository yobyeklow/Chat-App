package validation

import (
	"fmt"
	"strings"
	"web_socket/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() error {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("Failed to get validator engine")
	}
	RegisterCustomValidation(validate)
	return nil
}

func HandleValidationErrors(err error) gin.H {
	if valitionError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, e := range valitionError {
			root := strings.Split(e.Namespace(), ".")[0]
			rawPath := strings.TrimPrefix(e.Namespace(), root+".")
			parts := strings.Split(rawPath, ".")
			for i, part := range parts {
				if strings.Contains("parts", "[") {
					idx := strings.Index(part, "[")
					base := utils.CamelToSnake(part[:idx])
					index := part[idx:]
					parts[i] = base + index
				} else {
					parts[i] = utils.CamelToSnake(part)
				}
			}
			fieldPath := strings.Join(parts, ".")

			switch e.Tag() {
			case "gt":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("%s must be less than %s", fieldPath, e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than or equal to %s", fieldPath, e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("%s must be less than or equal to %s", fieldPath, e.Param())
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid UUID", fieldPath)
			case "short_key":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid ShortKey", fieldPath)
			case "slug":
				errors[fieldPath] = fmt.Sprintf("%s can only contain lowercase letters, numbers, hyphens, or periods.", fieldPath)
			case "min":
				errors[fieldPath] = fmt.Sprintf("%s must be more than %s characters", fieldPath, e.Param())
			case "max":
				errors[fieldPath] = fmt.Sprintf("%s must be less than %s ký tự", fieldPath, e.Param())
			case "min_int":
				errors[fieldPath] = fmt.Sprintf("%s must have a greater value %s", fieldPath, e.Param())
			case "max_int":
				errors[fieldPath] = fmt.Sprintf("%s must have a smaller value %s", fieldPath, e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[fieldPath] = fmt.Sprintf("%s must be one of the values: %s", fieldPath, allowedValues)
			case "required":
				errors[fieldPath] = fmt.Sprintf("%s is required", fieldPath)
			case "search":
				errors[fieldPath] = fmt.Sprintf("%s can only contain lowercase letters, uppercase, numbers và blank", fieldPath)
			case "email":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid Email", fieldPath)
			case "datetime":
				errors[fieldPath] = fmt.Sprintf("%s should be format YYYY-MM-DD", fieldPath)
			case "email_advanced":
				errors[fieldPath] = fmt.Sprintf("%s is on the email blacklist", fieldPath)
			case "password_strong":
				errors[fieldPath] = fmt.Sprintf("%s must contain at least 8 characters, including (lowercase letters, uppercase letters, numbers, and special characters)", fieldPath)
			case "file_ext":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[fieldPath] = fmt.Sprintf("%s only allow these files have extension: %s", fieldPath, allowedValues)
		}
		return gin.H{"error": errors}
	}
	return gin.H{
		"error":  "Something is wrong!",
		"detail": err.Error(),
	}
}
