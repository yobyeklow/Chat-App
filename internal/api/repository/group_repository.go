package repository

import (
	"context"
	"web_socket/internal/common/database/sqlc"

	"github.com/google/uuid"
)

type SQLGroupRepository struct {
	db sqlc.Querier
}

func NewSQLGroupRepository(db sqlc.Querier) GroupRepository {
	return &SQLGroupRepository{
		db: db,
	}
}

func (gr *SQLGroupRepository) CreateGroup(ctx context.Context, groupName string) (sqlc.Group, error) {
	groupData, err := gr.db.CreateGroup(ctx, groupName)
	if err != nil {
		return sqlc.Group{}, err
	}
	return groupData, nil
}
func (gr *SQLGroupRepository) GetAllGroups(ctx context.Context, arg sqlc.GetAllGroupsParams) ([]sqlc.GetAllGroupsRow, error) {
	groupData, err := gr.db.GetAllGroups(ctx, arg)
	if err != nil {
		return []sqlc.GetAllGroupsRow{}, err
	}
	return groupData, nil
}
func (gr *SQLGroupRepository) UpdateGroup(ctx context.Context, arg sqlc.UpdateGroupParams) (sqlc.Group, error) {
	groupData, err := gr.db.UpdateGroup(ctx, arg)
	if err != nil {
		return sqlc.Group{}, err
	}
	return groupData, nil
}
func (gr *SQLGroupRepository) SoftDeleteGroup(ctx context.Context, groupUUID uuid.UUID) (sqlc.Group, error) {
	groupData, err := gr.db.SoftDeleteGroup(ctx, groupUUID)
	if err != nil {
		return sqlc.Group{}, err
	}
	return groupData, nil
}
func (gr *SQLGroupRepository) HardDeleteGroup(ctx context.Context, groupUuid uuid.UUID) error {
	err := gr.db.HardDeleteGroup(ctx, groupUuid)
	if err != nil {
		return err
	}
	return nil
}
func (gr *SQLGroupRepository) LeaveGroup(ctx context.Context, arg sqlc.LeaveGroupParams) error {
	err := gr.db.LeaveGroup(ctx, arg)
	if err != nil {
		return err
	}
	return nil
}
func (gr *SQLGroupRepository) AddMemberToGroup(ctx context.Context, arg sqlc.AddMemberToGroupParams) (sqlc.GroupMember, error) {
	groupData, err := gr.db.AddMemberToGroup(ctx, arg)
	if err != nil {
		return sqlc.GroupMember{}, err
	}
	return groupData, nil
}
func (gr *SQLGroupRepository) GetGroupMembers(ctx context.Context, arg sqlc.GetGroupMembersParams) ([]sqlc.GetGroupMembersRow, error) {
	members, err := gr.db.GetGroupMembers(ctx, arg)
	if err != nil {
		return []sqlc.GetGroupMembersRow{}, err
	}
	return members, nil
}
func (gr *SQLGroupRepository) GetGroupMemberRole(ctx context.Context, arg sqlc.GetMemberRoleParams) (int32, error) {
	role, err := gr.db.GetMemberRole(ctx, arg)
	if err != nil {
		return 0, err
	}
	return role, nil
}
func (gr *SQLGroupRepository) GetMemberInfo(ctx context.Context, arg sqlc.GetMemberInfoParams) (sqlc.GetMemberInfoRow, error) {
	memberData, err := gr.db.GetMemberInfo(ctx, arg)
	if err != nil {
		return sqlc.GetMemberInfoRow{}, err
	}
	return memberData, nil
}
func (gr *SQLGroupRepository) UpdateMemberRole(ctx context.Context, arg sqlc.UpdateMemberRoleParams) (sqlc.GroupMember, error) {
	member, err := gr.db.UpdateMemberRole(ctx, arg)
	if err != nil {
		return sqlc.GroupMember{}, err
	}
	return member, nil
}
func (gr *SQLGroupRepository) RemoveMember(ctx context.Context, arg sqlc.RemoveMemberParams) (sqlc.GroupMember, error) {
	member, err := gr.db.RemoveMember(ctx, arg)
	if err != nil {
		return sqlc.GroupMember{}, err
	}
	return member, nil
}
