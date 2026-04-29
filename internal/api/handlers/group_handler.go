package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"web_socket/internal/api/dto"
	"web_socket/internal/api/services"
	"web_socket/internal/common/utils"
	"web_socket/internal/common/validation"
	"web_socket/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupHandler struct {
	service      services.GroupService
	tokenService auth.TokenService
}

func NewGroupHandler(service services.GroupService, tokenService auth.TokenService) *GroupHandler {
	return &GroupHandler{
		service:      service,
		tokenService: tokenService,
	}
}

func (gh *GroupHandler) CreateGroup(ctx *gin.Context) {
	var input dto.GroupInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUUIDData, _ := ctx.Get("user_uuid")
	userUUID, err := uuid.Parse(userUUIDData.(string))
	groupData, err := gh.service.CreateGroup(ctx, userUUID, input.GroupName)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	groupDTO := dto.MapToGroupDTO(groupData)
	utils.ResponseSuccess(ctx, http.StatusOK, "Created group successfully!", groupDTO)
}
func (gh *GroupHandler) GetAllGroups(ctx *gin.Context) {
	var input dto.GroupSearchParams
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	authHeader := ctx.GetHeader("Authorization")
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	payload, err := gh.tokenService.DecryptAccessTokenPayload(accessToken)
	if err != nil {
		utils.ResponseError(ctx, err)
		return

	}
	fmt.Println(payload)
}
