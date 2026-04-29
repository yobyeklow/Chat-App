package handlers

import (
	"net/http"
	"time"
	"web_socket/internal/api/dto"
	"web_socket/internal/api/services"
	"web_socket/internal/common/utils"
	"web_socket/internal/common/validation"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	services services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{
		services: service,
	}
}

func (ah *AuthHandler) Register(ctx *gin.Context) {
	var input dto.UserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userInput := input.MapCreateInputToModel()
	userData, err := ah.services.Register(ctx, userInput)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapToUserDTO(userData)
	utils.ResponseSuccess(ctx, http.StatusOK, "Created user successfully!", userDTO)
}
func (ah *AuthHandler) Login(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	accessToken, refreshToken, expiredAt, err := ah.services.Login(ctx, input.Email, input.Password)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	resp := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(time.Duration(expiredAt) * time.Second).Format("2006-01-02 15:04:05"),
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Login successfully!", resp)
}
func (ah *AuthHandler) Logout(ctx *gin.Context) {
	var input dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	err := ah.services.Logout(ctx, input.RefreshToken)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Logout successfully!")
}
func (ah *AuthHandler) RefreshToken(ctx *gin.Context) {
	var input dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	accessToken, refreshToken, expiredAt, err := ah.services.RefreshToken(ctx, input.RefreshToken)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	resp := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(time.Duration(expiredAt) * time.Second).Format("2006-01-02 15:04:05"),
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Reset token successfully!", resp)
}
