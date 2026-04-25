package handlers

import (
	"net/http"
	"web_socket/internal/dto"
	"web_socket/internal/services"
	"web_socket/internal/utils"
	"web_socket/internal/validation"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input dto.UserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userInput := input.MapCreateInputToModel()
	userData, err := uh.service.CreateUser(ctx, userInput)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapToUserDTO(userData)
	utils.ResponseSuccess(ctx, http.StatusOK, "Created user successfully!", userDTO)
}
func (uh *UserHandler) FindUserByEmail(ctx *gin.Context) {
	var input dto.GetUserByEmailParam
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	input.Email = utils.NormalizeString(input.Email)
	userData, err := uh.service.FindUserByEmail(ctx, input.Email)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapToUserDTO(userData)
	utils.ResponseSuccess(ctx, http.StatusOK, "User found successfully!", userDTO)
}
