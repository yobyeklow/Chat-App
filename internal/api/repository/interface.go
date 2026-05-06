package repository

import (
	"context"
	"web_socket/internal/common/database/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	FindUserByEmail(ctx context.Context, userEmail string) (sqlc.User, error)
	FindUserByUUID(ctx context.Context, userUUID uuid.UUID) (sqlc.User, error)
	SoftDeleteUser(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
	HardDeleteUser(ctx context.Context, userUuid uuid.UUID) error
	RestoreUser(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
}

type GroupRepository interface {
	CreateGroup(ctx context.Context, arg sqlc.CreateGroupParams) (sqlc.CreateGroupRow, error)
	GetAllGroups(ctx context.Context, arg sqlc.GetAllGroupsParams) ([]sqlc.GetAllGroupsRow, error)
	UpdateGroup(ctx context.Context, arg sqlc.UpdateGroupParams) (sqlc.Group, error)
	SoftDeleteGroup(ctx context.Context, groupUuid uuid.UUID) (sqlc.Group, error)
	HardDeleteGroup(ctx context.Context, groupUuid uuid.UUID) error
	CountGroups(ctx context.Context, search string, deleted bool) (int64, error)

	LeaveGroup(ctx context.Context, arg sqlc.LeaveGroupParams) error
	AddMemberToGroup(ctx context.Context, arg sqlc.AddMemberToGroupParams) error
	GetGroupMembers(ctx context.Context, arg sqlc.GetGroupMembersParams) ([]sqlc.GetGroupMembersRow, error)
	GetGroupMemberRole(ctx context.Context, arg sqlc.GetMemberRoleParams) (int32, error)
	GetMemberInfo(ctx context.Context, arg sqlc.GetMemberInfoParams) (sqlc.GetMemberInfoRow, error)
	UpdateMemberRole(ctx context.Context, arg sqlc.UpdateMemberRoleParams) (sqlc.GroupMember, error)
	RemoveMember(ctx context.Context, arg sqlc.RemoveMemberParams) (sqlc.GroupMember, error)
}
type ChatRepository interface {
	CreateDMConversation(ctx context.Context) (int32, error)
	FindDMConversationn(ctx context.Context, arg sqlc.FindDMConversationParams) (int32, error)
	AddParticipantToConversation(ctx context.Context, arg sqlc.AddParticipantToConversationParams) error
	SendMessage(ctx context.Context, arg sqlc.SendMessageParams) (int32, error)
	GetMessage(ctx context.Context, arg sqlc.GetMessagesParams) ([]sqlc.GetMessagesRow, error)
	ValidateReplyMessage(ctx context.Context, arg sqlc.ValidateReplyParams) (bool, error)
	GetDMConversations(ctx context.Context, curUserID int32) ([]sqlc.GetDMConversationRow, error)
	GetGroupConversations(ctx context.Context, curUserID int32) ([]sqlc.GetGroupConversationRow, error)
}
