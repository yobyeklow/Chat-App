package dto

import (
	"time"
	"web_socket/internal/common/database/sqlc"

	"github.com/google/uuid"
)

type GroupCreateInput struct {
	GroupName string `json:"group_name" binding:"min=3,max=50,search,required"`
}
type AddMemberInput struct {
	MemberRole int32  `json:"member_role" binding:"required,oneof=1 2 3"`
	UserUuid   string `json:"user_uuid" binding:"required"`
}
type GroupInputJSON struct {
	GroupName  string `json:"group_name" binding:"omitempty,min=3,max=50,search"`
	MemberRole int32  `json:"member_role" binding:"omitempty,oneof=1 2 3"`
}
type GroupInputURI struct {
	GroupUuid string `uri:"uuid" binding:"uuid"`
	UserUuid  string `uri:"user_uuid" binding:"omitempty,uuid"`
}

type GroupDTO struct {
	GroupUUID string `json:"group_uuid"`
	GroupName string `json:"group_name"`
	CreatedAt string `json:"CreatedAt"`
	DeletedAt string `json:"DeletedAt"`
}
type GetAllGroupsRow struct {
	GroupUuid      uuid.UUID `json:"group_uuid"`
	GroupName      string    `json:"group_name"`
	GroupCreatedAt time.Time `json:"group_created_at"`
	GroupUpdatedAt time.Time `json:"group_updated_at"`
	MemberRole     int32     `json:"member_role"`
	JointedAt      time.Time `json:"jointed_at"`
}
type GroupSearchParams struct {
	Search  string `form:"search" binding:"omitempty,min=3,max=50,search"`
	Page    int32  `form:"page" binding:"omitempty,gte=1"`
	Limit   int32  `form:"limit" binding:"omitempty,gte=1,lte=500"`
	Deleted bool   `form:"deleted" binding:"omitempty"`
}
type GroupMemberDTO struct {
	MemberRole string `json:"member_role"`
	JointedAt  string `json:"jointed_at"`
	UserUUID   string `json:"user_uuid"`
	GroupUUID  string `json:"group_uuid"`
}

func MapToGroupDTO(groupInput sqlc.Group) *GroupDTO {
	dto := &GroupDTO{
		GroupUUID: groupInput.GroupUuid.String(),
		GroupName: groupInput.GroupName,
		CreatedAt: groupInput.GroupCreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if groupInput.GroupDeletedAt.Valid {
		dto.DeletedAt = groupInput.GroupDeletedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	} else {
		dto.DeletedAt = ""
	}
	return dto
}
func MapToGroupMemberDTO(groupInput sqlc.GroupMember, userUUID uuid.UUID, groupUUID uuid.UUID) *GroupMemberDTO {
	dto := &GroupMemberDTO{
		UserUUID:   userUUID.String(),
		GroupUUID:  groupUUID.String(),
		MemberRole: mapMemberRoleToText(groupInput.MemberRole),
		JointedAt:  groupInput.JointedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return dto
}
func MapGroupMemberssToDTO(groupsData []sqlc.GetGroupMembersRow) []GroupMemberDTO {
	dtos := make([]GroupMemberDTO, 0, len(groupsData))
	for _, group := range groupsData {
		dto := &GroupMemberDTO{
			GroupUUID:  group.GroupUuid.String(),
			UserUUID:   group.UserUuid.String(),
			MemberRole: mapMemberRoleToText(group.MemberRole),
			JointedAt:  group.JointedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		dtos = append(dtos, *dto)
	}
	return dtos
}
func MapGroupsToDTO(groupsData []sqlc.GetAllGroupsRow) []GroupDTO {
	dtos := make([]GroupDTO, 0, len(groupsData))
	for _, group := range groupsData {
		dto := &GroupDTO{
			GroupUUID: group.GroupUuid.String(),
			GroupName: group.GroupName,
			CreatedAt: group.GroupCreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		dtos = append(dtos, *dto)
	}
	return dtos
}
func mapMemberRoleToText(memberRole int32) string {
	switch memberRole {
	case 1:
		return "Member"
	case 2:
		return "Moderator"
	case 3:
		return "Admin"
	default:
		return "None"
	}
}
