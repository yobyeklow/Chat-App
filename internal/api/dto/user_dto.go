package dto

import "web_socket/internal/common/database/sqlc"

type UserInput struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email" binding:"required,email,email_advanced"`
	Password string `json:"password" binding:"required"`
	Status   int32  `json:"status" binding:"omitempty,oneof=1 2 3"`
	Role     int32  `json:"role" binding:"omitempty,oneof=1 2"`
}
type UserDTO struct {
	Email     string `json:"email" binding:"required,email,email_advanced"`
	Status    string `json:"status" binding:"omitempty,oneof=1 2 3"`
	Role      string `json:"role" binding:"omitempty,oneof=1 2"`
	CreatedAt string `json:"CreatedAt"`
	DeletedAt string `json:"DeletedAt"`
}
type GetUserByEmailParam struct {
	Email string `uri:"email" binding:"email,email_advanced"`
}
type GetUserByUUIDParam struct {
	Uuid string `uri:"uuid" binding:"uuid"`
}

func (input *UserInput) MapCreateInputToModel() sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		UserEmail:    input.Email,
		UserPassword: input.Password,
		UserStatus:   input.Status,
		UserRole:     input.Role,
	}
}
func MapToUserDTO(userInput sqlc.User) *UserDTO {
	dto := &UserDTO{
		Email:     userInput.UserEmail,
		Status:    mapStatusText(int(userInput.UserStatus)),
		Role:      mapStatusText(int(userInput.UserRole)),
		CreatedAt: userInput.UserCreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if userInput.UserDeletedAt.Valid {
		dto.DeletedAt = userInput.UserDeletedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	} else {
		dto.DeletedAt = ""
	}
	return dto
}
func mapStatusText(status int) string {
	switch status {
	case 1:
		return "Active"
	case 2:
		return "Inactive"
	case 3:
		return "Banned"
	default:
		return "None"
	}
}
func mapRoleText(status int) string {
	switch status {
	case 1:
		return "User"
	case 2:
		return "Adminstrator"
	default:
		return "None"
	}
}
