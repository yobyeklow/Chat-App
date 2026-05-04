package services

import (
	"context"
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
	OwnerGroupRole   int32 = 3
	MemberRoleGr     int32 = 1
	ModeratorRoleGr  int32 = 2
	AdminAccountRole int32 = 2
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
		MemberRole: OwnerGroupRole,
	}
	_, err = gs.repo.AddMemberToGroup(context, arg)
	if err != nil {
		return sqlc.Group{}, utils.WrapError("Failed to create group", utils.ErrCodeInternal, err)
	}
	return groupData, nil
}
func (gs *groupService) GetAllGroups(ctx *gin.Context, userUUID uuid.UUID, search string, page int32, limit int32, deleted bool) ([]sqlc.GetAllGroupsRow, int32, error) {
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
		Search:    search,
	}
	groupsData, err := gs.repo.GetAllGroups(context, arg)
	if err != nil {
		return []sqlc.GetAllGroupsRow{}, 0, utils.WrapError("Failed to get all groups", utils.ErrCodeInternal, err)
	}
	total, err := gs.repo.CountGroups(context, search, deleted)
	if err != nil {
		return []sqlc.GetAllGroupsRow{}, 0, utils.WrapError("Failed to count groups", utils.ErrCodeInternal, err)
	}
	return groupsData, int32(total), nil
}
func (gs *groupService) UpdateGroup(ctx *gin.Context, userUUID uuid.UUID, userRole int32, groupName string, groupUUID uuid.UUID) (sqlc.Group, error) {
	context := ctx.Request.Context()
	if strings.TrimSpace(groupName) == "" {
		return sqlc.Group{}, utils.NewError("Group name is required", utils.ErrCodeBadRequest)
	}
	err := gs.verifyOwnerOrAdmin(context, userUUID, userRole, groupUUID)
	if err != nil {
		return sqlc.Group{}, err
	}
	updateGroupArg := sqlc.UpdateGroupParams{
		GroupName: groupName,
		GroupUuid: groupUUID,
	}
	groupData, err := gs.repo.UpdateGroup(context, updateGroupArg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Group{}, utils.NewError("Group does not exist", utils.ErrCodeNotFound)
		}
		return sqlc.Group{}, utils.WrapError("Failed to update group", utils.ErrCodeInternal, err)
	}
	return groupData, nil
}
func (gs *groupService) SoftDeleteGroup(ctx *gin.Context, userRole int32, userUUID uuid.UUID, groupUUID uuid.UUID) (sqlc.Group, error) {
	context := ctx.Request.Context()
	err := gs.verifyOwnerOrAdmin(context, userUUID, userRole, groupUUID)
	if err != nil {
		return sqlc.Group{}, err
	}
	groupData, err := gs.repo.SoftDeleteGroup(context, groupUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Group{}, utils.NewError("Group does not exist", utils.ErrCodeNotFound)
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
			return utils.NewError("Group does not exist", utils.ErrCodeNotFound)
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
		return utils.NewError("User does not belong to this group.", utils.ErrCodeBadRequest)
	}
	if role == OwnerGroupRole {
		return utils.NewError("You should transfer 'Owner Key' before do this action!", utils.ErrCodeUnauthorized)
	}

	leaveGrArg := sqlc.LeaveGroupParams{
		UserUuid:  userUUID,
		GroupUuid: groupUUID,
	}
	err = gs.repo.LeaveGroup(context, leaveGrArg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("Group does not exist", utils.ErrCodeNotFound)
		}
		return utils.WrapError("Failed to leave group", utils.ErrCodeInternal, err)
	}
	return nil
}
func (gs *groupService) JoinGroup(ctx *gin.Context, groupUUID uuid.UUID, userUUID uuid.UUID, memberRole int32) error {
	context := ctx.Request.Context()
	roleArg := sqlc.GetMemberRoleParams{
		UserUuid:  userUUID,
		GroupUuid: groupUUID,
	}
	_, err := gs.repo.GetGroupMemberRole(context, roleArg)
	if err == nil {
		return utils.NewError("This user has joined the group", utils.ErrCodeBadRequest)
	}

	joinGrArg := sqlc.AddMemberToGroupParams{
		GroupUuid:  groupUUID,
		UserUuid:   userUUID,
		MemberRole: memberRole,
	}
	_, err = gs.repo.AddMemberToGroup(context, joinGrArg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("Group does not existe", utils.ErrCodeBadRequest)
		}
		return utils.WrapError("Failed to join group", utils.ErrCodeInternal, err)
	}
	return nil
}
func (gs *groupService) verifyOwnerOrAdmin(ctx context.Context, userUUID uuid.UUID, userRole int32, groupUUID uuid.UUID) error {
	if userRole != AdminAccountRole {
		arg := sqlc.GetMemberRoleParams{
			UserUuid:  userUUID,
			GroupUuid: groupUUID,
		}
		role, err := gs.repo.GetGroupMemberRole(ctx, arg)
		if err != nil {
			return utils.NewError("User does not belong to this group.", utils.ErrCodeBadRequest)
		}
		if role != OwnerGroupRole {
			return utils.NewError("You're not allowed to do this!", utils.ErrCodeUnauthorized)
		}
	}
	return nil
}
func (gs *groupService) GetGroupMembers(ctx *gin.Context, groupUUID uuid.UUID, userUUID uuid.UUID, page int32, limit int32) ([]sqlc.GetGroupMembersRow, error) {
	context := ctx.Request.Context()
	roleCheckArg := sqlc.GetMemberRoleParams{
		UserUuid:  userUUID,
		GroupUuid: groupUUID,
	}
	_, err := gs.repo.GetGroupMemberRole(ctx, roleCheckArg)
	if err != nil {
		return []sqlc.GetGroupMembersRow{}, utils.NewError("User does not belong to this group.", utils.ErrCodeBadRequest)
	}
	if limit <= 0 {
		envLimit := utils.GetEnvInt("LIMIT_ITEM_ON_PER_PAGE", 10)
		limit = int32(envLimit)
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	getGrArg := sqlc.GetGroupMembersParams{
		GroupUuid: groupUUID,
		Offsetarg: offset,
		Limitarg:  limit,
	}
	groupsData, err := gs.repo.GetGroupMembers(context, getGrArg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []sqlc.GetGroupMembersRow{}, utils.NewError(
				"Group does not exist", utils.ErrCodeBadRequest,
			)
		}
		return []sqlc.GetGroupMembersRow{}, utils.WrapError("Failed to fetch group members", utils.ErrCodeInternal, err)
	}
	return groupsData, nil
}
func (gs *groupService) UpdateMemberRole(ctx *gin.Context, memberRole int32, groupUUID uuid.UUID, userUUID uuid.UUID) (sqlc.GroupMember, error) {
	context := ctx.Request.Context()

	userUUIDData := ctx.GetString("user_uuid")
	curUserUUID, err := uuid.Parse(userUUIDData)
	if err != nil {
		return sqlc.GroupMember{}, utils.NewError("Invalid User UUID request", utils.ErrCodeBadRequest)
	}
	curUserRole := ctx.GetInt32("user_role")
	err = gs.verifyOwnerOrAdmin(context, curUserUUID, curUserRole, groupUUID)
	if err != nil {
		return sqlc.GroupMember{}, err
	}
	arg := sqlc.UpdateMemberRoleParams{
		MemberRole: memberRole,
		GroupUuid:  groupUUID,
		UserUuid:   userUUID,
	}
	userData, err := gs.repo.UpdateMemberRole(context, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.GroupMember{}, utils.NewError(
				"User is not a member of this group", utils.ErrCodeBadRequest,
			)
		}
		return sqlc.GroupMember{}, utils.WrapError("Failed to update member role", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
func (gs *groupService) RemoveMember(ctx *gin.Context, groupUUID uuid.UUID, userUUID uuid.UUID) (sqlc.GroupMember, error) {
	context := ctx.Request.Context()
	userUUIDData := ctx.GetString("user_uuid")
	curUserUUID, err := uuid.Parse(userUUIDData)
	if err != nil {
		return sqlc.GroupMember{}, utils.NewError("Invalid User UUID request", utils.ErrCodeBadRequest)
	}
	curUserRole := ctx.GetInt32("user_role")
	err = gs.verifyOwnerOrAdmin(context, curUserUUID, curUserRole, groupUUID)
	if err != nil {
		return sqlc.GroupMember{}, err
	}
	arg := sqlc.RemoveMemberParams{
		GroupUuid: groupUUID,
		UserUuid:  userUUID,
	}
	userData, err := gs.repo.RemoveMember(context, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.GroupMember{}, utils.NewError(
				"User is not a member of this group", utils.ErrCodeBadRequest,
			)
		}
		return sqlc.GroupMember{}, utils.WrapError("Failed to update member role", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
func (gs *groupService) GetMemberInfo(ctx *gin.Context, groupUUID uuid.UUID, curUserUUID uuid.UUID, targetUserUUID uuid.UUID) (sqlc.GetMemberInfoRow, error) {
	context := ctx.Request.Context()
	roleCheckArg := sqlc.GetMemberRoleParams{
		UserUuid:  curUserUUID,
		GroupUuid: groupUUID,
	}
	_, err := gs.repo.GetGroupMemberRole(ctx, roleCheckArg)
	if err != nil {
		return sqlc.GetMemberInfoRow{}, utils.NewError("User does not belong to this group.", utils.ErrCodeBadRequest)
	}
	arg := sqlc.GetMemberInfoParams{
		GroupUuid: groupUUID,
		UserUuid:  targetUserUUID,
	}
	userData, err := gs.repo.GetMemberInfo(context, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.GetMemberInfoRow{}, utils.NewError(
				"User is not a member of this group", utils.ErrCodeBadRequest,
			)
		}
		return sqlc.GetMemberInfoRow{}, utils.WrapError("Failed to update member role", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
