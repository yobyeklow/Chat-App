package handlers

import (
	"net/http"
	"web_socket/internal/dto"
	"web_socket/internal/services"
	"web_socket/internal/utils"
	"web_socket/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
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
func (uh *UserHandler) FindUserByUUID(ctx *gin.Context) {
	var input dto.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUID, err := uuid.Parse(input.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userData, err := uh.service.FindUserByUUID(ctx, userUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapToUserDTO(userData)
	utils.ResponseSuccess(ctx, http.StatusOK, "User found successfully!", userDTO)
}
func (uh *UserHandler) SoftDeleteUser(ctx *gin.Context) {
	var input dto.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUID, err := uuid.Parse(input.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userData, err := uh.service.SoftDeleteUser(ctx, userUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapToUserDTO(userData)
	utils.ResponseSuccess(ctx, http.StatusOK, "Deleted user successfully!", userDTO)
}
func (uh *UserHandler) HardDeleteUser(ctx *gin.Context) {
	var input dto.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUID, err := uuid.Parse(input.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	err = uh.service.HardDeleteUser(ctx, userUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}
func (uh *UserHandler) RestoreUser(ctx *gin.Context) {
	var input dto.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUID, err := uuid.Parse(input.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userData, err := uh.service.RestoreUser(ctx, userUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapToUserDTO(userData)
	utils.ResponseSuccess(ctx, http.StatusOK, "Restored user successfully!", userDTO)
}
