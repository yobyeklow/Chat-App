package services

import (
	"encoding/json"
	"time"
	"web_socket/internal/common/database/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	FindUserByEmail(ctx *gin.Context, userEmail string) (sqlc.User, error)
	FindUserByUUID(ctx *gin.Context, userUUID uuid.UUID) (sqlc.User, error)
	SoftDeleteUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)
	HardDeleteUser(ctx *gin.Context, userUuid uuid.UUID) error
	RestoreUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)
}

type AuthService interface {
	Register(ctx *gin.Context, userInput sqlc.CreateUserParams) (sqlc.User, error)
	Login(ctx *gin.Context, email string, password string) (string, string, int, error)
	Logout(ctx *gin.Context, refreshTokenStr string) error
	RefreshToken(ctx *gin.Context, refreshTokenStr string) (string, string, int, error)
}
type GroupService interface {
	CreateGroup(ctx *gin.Context, groupName string, userUUID uuid.UUID) (sqlc.CreateGroupRow, error)
	GetAllGroups(ctx *gin.Context, userUUID uuid.UUID, search string, page int32, limit int32, deleted bool) ([]sqlc.GetAllGroupsRow, int32, error)
	UpdateGroup(ctx *gin.Context, userUUID uuid.UUID, userRole int32, groupName string, groupUUID uuid.UUID) (sqlc.Group, error)
	SoftDeleteGroup(ctx *gin.Context, userRole int32, userUUID uuid.UUID, groupUUID uuid.UUID) (sqlc.Group, error)
	HardDeleteGroup(ctx *gin.Context, groupUUID uuid.UUID) error
	LeaveGroup(ctx *gin.Context, userUUID uuid.UUID, groupUUID uuid.UUID) error
	JoinGroup(ctx *gin.Context, groupUUID uuid.UUID, userUUID uuid.UUID) error
	GetMemberRole(ctx *gin.Context, userUUID uuid.UUID, groupUUID uuid.UUID) (int32, error)
	GetGroupMembers(ctx *gin.Context, groupUUID uuid.UUID, userUUID uuid.UUID, page int32, limit int32) ([]sqlc.GetGroupMembersRow, error)
	UpdateMemberRole(ctx *gin.Context, memberRole int32, groupUUID uuid.UUID, userUUID uuid.UUID) (sqlc.GroupMember, error)
	RemoveMember(ctx *gin.Context, groupUUID uuid.UUID, userUUID uuid.UUID) (sqlc.GroupMember, error)
	GetMemberInfo(ctx *gin.Context, groupUUID uuid.UUID, curUserUUID uuid.UUID, targetUserUUID uuid.UUID) (sqlc.GetMemberInfoRow, error)
}
type ChatService interface {
	StartDirectMessage(ctx *gin.Context, sender uuid.UUID, receiver uuid.UUID) error
	SendMessage(ctx *gin.Context, senderUUID uuid.UUID, conversationID int32, content string, attachments json.RawMessage, msgType int32, replyTo *int32) (int32, error)
	GetConversations(ctx *gin.Context, curUserUUID uuid.UUID) ([]Conversation, error)
	GetMessages(ctx *gin.Context, conversationID int32, cursorTime *time.Time, cursorID *int32, limit int32) ([]sqlc.GetMessagesRow, error)
}
