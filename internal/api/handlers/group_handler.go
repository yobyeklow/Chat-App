package handlers

import (
	"net/http"
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
	var input dto.GroupCreateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUUIDData := ctx.GetString("user_uuid")
	userUUID, err := uuid.Parse(userUUIDData)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User ID", utils.ErrCodeBadRequest))
		return
	}
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
	if err := ctx.ShouldBindQuery(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUIDData := ctx.GetString("user_uuid")
	userUUID, err := uuid.Parse(userUUIDData)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User ID", utils.ErrCodeBadRequest))
		return
	}
	groupsData, totalRecords, err := gh.service.GetAllGroups(ctx, userUUID, input.Search, input.Page, input.Limit, input.Deleted)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	groupsDTO := dto.MapGroupsToDTO(groupsData)
	paginationResp := utils.NewPaginationResponse(groupsDTO, input.Page, input.Limit, totalRecords)
	utils.ResponseSuccess(ctx, http.StatusOK, "Fetched groups successfully!", paginationResp)

}
func (gh *GroupHandler) UpdateGroup(ctx *gin.Context) {
	var inputURI dto.GroupInputURI
	if err := ctx.ShouldBindUri(&inputURI); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	var inputJSON dto.GroupInputJSON
	if err := ctx.ShouldBindJSON(&inputJSON); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUIDData := ctx.GetString("user_uuid")
	userUUID, err := uuid.Parse(userUUIDData)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User ID", utils.ErrCodeBadRequest))
		return
	}
	userRole := ctx.GetInt32("user_role")
	groupUUID, err := uuid.Parse(inputURI.GroupUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid Group ID", utils.ErrCodeBadRequest))
		return
	}
	groupData, err := gh.service.UpdateGroup(ctx, userUUID, userRole, inputJSON.GroupName, groupUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	groupDTO := dto.MapToGroupDTO(groupData)
	utils.ResponseSuccess(ctx, http.StatusOK, "Updated group successfully", groupDTO)
}
func (gh *GroupHandler) SoftDeleteGroup(ctx *gin.Context) {
	var input dto.GroupInputURI
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUIDData := ctx.GetString("user_uuid")
	userUUID, err := uuid.Parse(userUUIDData)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User ID", utils.ErrCodeBadRequest))
		return
	}
	userRoleData := ctx.GetInt32("user_role")
	groupUUID, err := uuid.Parse(input.GroupUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid Group ID", utils.ErrCodeBadRequest))
		return
	}
	groupData, err := gh.service.SoftDeleteGroup(ctx, userRoleData, userUUID, groupUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	groupDTO := dto.MapToGroupDTO(groupData)
	utils.ResponseSuccess(ctx, http.StatusOK, "Deleted group successfully", groupDTO)
}
func (gh *GroupHandler) LeaveGroup(ctx *gin.Context) {
	var input dto.GroupInputURI
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUIDData := ctx.GetString("user_uuid")
	userUUID, err := uuid.Parse(userUUIDData)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User ID", utils.ErrCodeBadRequest))
		return
	}

	groupUUID, err := uuid.Parse(input.GroupUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid Group ID", utils.ErrCodeBadRequest))
		return
	}
	err = gh.service.LeaveGroup(ctx, userUUID, groupUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Successfully left the group")
}
func (gh *GroupHandler) AddMemberToGroup(ctx *gin.Context) {
	var inputURI dto.GroupInputURI
	if err := ctx.ShouldBindUri(&inputURI); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	var inputJSON dto.AddMemberInput
	if err := ctx.ShouldBindJSON(&inputJSON); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUUID, err := uuid.Parse(inputJSON.UserUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User ID", utils.ErrCodeBadRequest))
		return
	}
	groupUUID, err := uuid.Parse(inputURI.GroupUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid Group ID", utils.ErrCodeBadRequest))
		return
	}
	err = gh.service.JoinGroup(ctx, groupUUID, userUUID, inputJSON.MemberRole)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Added member to group")
}
func (gh *GroupHandler) HardDeleteGroup(ctx *gin.Context) {
	var input dto.GroupInputURI
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	groupUUID, err := uuid.Parse(input.GroupUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid Group ID", utils.ErrCodeBadRequest))
		return
	}
	err = gh.service.HardDeleteGroup(ctx, groupUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusOK)
}
func (gh *GroupHandler) GetGroupMembers(ctx *gin.Context) {
	var inputUri dto.GroupInputURI
	if err := ctx.ShouldBindUri(&inputUri); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	var inputQuery dto.GroupSearchParams
	if err := ctx.ShouldBindQuery(&inputQuery); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUIDData := ctx.GetString("user_uuid")
	userUUID, err := uuid.Parse(userUUIDData)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User ID", utils.ErrCodeBadRequest))
		return
	}
	groupUUID, err := uuid.Parse(inputUri.GroupUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid Group ID", utils.ErrCodeBadRequest))
		return
	}
	groupMembers, err := gh.service.GetGroupMembers(ctx, groupUUID, userUUID, inputQuery.Page, inputQuery.Limit)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	groupMembersDTO := dto.MapGroupMemberssToDTO(groupMembers)
	utils.ResponseSuccess(ctx, http.StatusOK, "Fetched group members successfully", groupMembersDTO)
}
func (gh *GroupHandler) GetMemberInfo(ctx *gin.Context) {
	var inputUri dto.GroupInputURI
	if err := ctx.ShouldBindUri(&inputUri); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	curUserUUIDData := ctx.GetString("user_uuid")
	curUserUUID, err := uuid.Parse(curUserUUIDData)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User ID", utils.ErrCodeBadRequest))
		return
	}
	groupUUID, err := uuid.Parse(inputUri.GroupUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid Group ID", utils.ErrCodeBadRequest))
		return
	}
	targetUserUUID, err := uuid.Parse(inputUri.UserUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid Target User ID", utils.ErrCodeBadRequest))
		return
	}
	memberData, err := gh.service.GetMemberInfo(ctx, groupUUID, curUserUUID, targetUserUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Fetched group members successfully", memberData)
}
func (gh *GroupHandler) UpdateMemberRole(ctx *gin.Context) {
	var inputURI dto.GroupInputURI
	if err := ctx.ShouldBindUri(&inputURI); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	var inputJSON dto.GroupInputJSON
	if err := ctx.ShouldBindJSON(&inputJSON); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUID, err := uuid.Parse(inputURI.UserUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User ID", utils.ErrCodeBadRequest))
		return
	}
	groupUUID, err := uuid.Parse(inputURI.GroupUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid Group ID", utils.ErrCodeBadRequest))
		return
	}
	groupMemberData, err := gh.service.UpdateMemberRole(ctx, inputJSON.MemberRole, groupUUID, userUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	groupMemberDTO := dto.MapToGroupMemberDTO(groupMemberData, userUUID, groupUUID)
	utils.ResponseSuccess(ctx, http.StatusOK, "Updated data successfully", groupMemberDTO)
}
func (gh *GroupHandler) RemoveMember(ctx *gin.Context) {
	var inputURI dto.GroupInputURI
	if err := ctx.ShouldBindUri(&inputURI); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUID, err := uuid.Parse(inputURI.UserUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User ID", utils.ErrCodeBadRequest))
		return
	}
	groupUUID, err := uuid.Parse(inputURI.GroupUuid)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid Group ID", utils.ErrCodeBadRequest))
		return
	}
	_, err = gh.service.RemoveMember(ctx, groupUUID, userUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Removed member successfully", nil)
}
