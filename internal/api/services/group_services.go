package services

import (
	"database/sql"
	"errors"
	"strings"
	"web_socket/internal/api/repository"
	"web_socket/internal/common/database/sqlc"
	"web_socket/internal/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type groupService struct {
	repo repository.GroupRepository
}

var (
	OwnerGroupRole   = 3
	AdminAccountRole = 2
)

func NewGroupService(repo repository.GroupRepository) GroupService {
	return &groupService{
		repo: repo,
	}
}
func (gs *groupService) CreateGroup(ctx *gin.Context, userUUID uuid.UUID, groupName string) (sqlc.Group, error) {
	context := ctx.Request.Context()
	if strings.TrimSpace(groupName) == "" {
		return sqlc.Group{}, utils.NewError("Group name is required", utils.ErrCodeBadRequest)
	}

	groupData, err := gs.repo.CreateGroup(context, groupName)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return sqlc.Group{}, utils.WrapError("Duplicate group name", utils.ErrCodeConflict, err)
		}
		return sqlc.Group{}, utils.WrapError("Failed to create group", utils.ErrCodeInternal, err)
	}
	arg := sqlc.AddMemberToGroupParams{
		GroupUuid:  groupData.GroupUuid,
		UserUuid:   userUUID,
		MemberRole: 3,
	}
	_, err = gs.repo.AddMemberToGroup(context, arg)
	if err != nil {
		return sqlc.Group{}, utils.WrapError("Failed to create group", utils.ErrCodeInternal, err)
	}
	return groupData, nil
}
func (gs *groupService) GetAllGroups(ctx *gin.Context, userUUID uuid.UUID, page int32, limit int32) ([]sqlc.GetAllGroupsRow, error) {
	context := ctx.Request.Context()
	if limit <= 0 {
		envLimit := utils.GetEnvInt("LIMIT_ITEM_ON_PER_PAGE", 10)
		limit = int32(envLimit)
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	arg := sqlc.GetAllGroupsParams{
		UserUuid:  userUUID,
		Limitarg:  limit,
		Offsetarg: offset,
	}
	groupsData, err := gs.repo.GetAllGroups(context, arg)
	if err != nil {
		return []sqlc.GetAllGroupsRow{}, utils.WrapError("Failed to get all groups", utils.ErrCodeInternal, err)
	}
	return groupsData, nil
}
func (gs *groupService) UpdateGroup(ctx *gin.Context, userUUID uuid.UUID, userRole int32, groupName string, groupUUID uuid.UUID) (sqlc.Group, error) {
	context := ctx.Request.Context()
	if strings.TrimSpace(groupName) == "" {
		return sqlc.Group{}, utils.NewError("Group name is required", utils.ErrCodeBadRequest)
	}

	if userRole != int32(AdminAccountRole) {
		memberRoleArg := sqlc.GetMemberRoleParams{
			UserUuid:  userUUID,
			GroupUuid: groupUUID,
		}
		role, err := gs.repo.GetGroupMemberRole(context, memberRoleArg)
		if err != nil {
			return sqlc.Group{}, utils.NewError("User not belong to this group", utils.ErrCodeBadRequest)
		}
		if role != int32(OwnerGroupRole) {
			return sqlc.Group{}, utils.NewError("You're not allowed to do this!", utils.ErrCodeUnauthorized)
		}
	}
	updateGroupArg := sqlc.UpdateGroupParams{
		GroupName: groupName,
		GroupUuid: groupUUID,
	}
	groupData, err := gs.repo.UpdateGroup(context, updateGroupArg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Group{}, utils.NewError("Group not existed", utils.ErrCodeUnauthorized)
		}
		return sqlc.Group{}, utils.WrapError("Failed to update group", utils.ErrCodeInternal, err)
	}
	return groupData, nil
}
func (gs *groupService) SoftDeleteGroup(ctx *gin.Context, userRole int32, userUUID uuid.UUID, groupUUID uuid.UUID) (sqlc.Group, error) {
	context := ctx.Request.Context()

	if userRole != int32(AdminAccountRole) {
		arg := sqlc.GetMemberRoleParams{
			UserUuid:  userUUID,
			GroupUuid: groupUUID,
		}
		role, err := gs.repo.GetGroupMemberRole(context, arg)
		if err != nil {
			return sqlc.Group{}, utils.NewError("User not belong to this group", utils.ErrCodeBadRequest)
		}
		if role != int32(OwnerGroupRole) {
			return sqlc.Group{}, utils.NewError("You're not allowed to do this!", utils.ErrCodeUnauthorized)
		}
	}

	groupData, err := gs.repo.SoftDeleteGroup(context, groupUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Group{}, utils.NewError("Group not existed!", utils.ErrCodeNotFound)
		}
		return sqlc.Group{}, utils.WrapError("Failed to delete group", utils.ErrCodeInternal, err)
	}
	return groupData, nil
}
func (gs *groupService) HardDeleteGroup(ctx *gin.Context, groupUUID uuid.UUID) error {
	context := ctx.Request.Context()
	err := gs.repo.HardDeleteGroup(context, groupUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("Group not existed!", utils.ErrCodeNotFound)
		}
		return utils.WrapError("Failed to delete group", utils.ErrCodeInternal, err)
	}
	return nil
}
func (gs *groupService) LeaveGroup(ctx *gin.Context, userUUID uuid.UUID, groupUUID uuid.UUID) error {
	context := ctx.Request.Context()
	roleArg := sqlc.GetMemberRoleParams{
		UserUuid:  userUUID,
		GroupUuid: groupUUID,
	}
	role, err := gs.repo.GetGroupMemberRole(context, roleArg)
	if err != nil {
		return utils.NewError("User not belong to this group", utils.ErrCodeBadRequest)
	}
	if role == int32(OwnerGroupRole) {
		return utils.NewError("You should transfer 'Owner Key' before do this action!", utils.ErrCodeUnauthorized)
	}

	leaveGrArg := sqlc.LeaveGroupParams{
		UserUuid:  userUUID,
		GroupUuid: groupUUID,
	}
	err = gs.repo.LeaveGroup(context, leaveGrArg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("Group not existed!", utils.ErrCodeNotFound)
		}
		return utils.WrapError("Failed to leave group", utils.ErrCodeInternal, err)
	}
	return nil
}
