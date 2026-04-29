package dto

import "web_socket/internal/common/database/sqlc"

type GroupInput struct {
	GroupName string `json:"group_name" binding:"required"`
}
type GroupDTO struct {
	GroupUUID string `json:"group_uuid"`
	GroupName string `json:"group_name"`
	CreatedAt string `json:"CreatedAt"`
	DeletedAt string `json:"DeletedAt"`
}
type GroupSearchParams struct {
	Search string `form:"search" binding:"omitempty,min=3,max=50,search"`
	Page   int32  `form:"page" binding:"omitempty,gte=1"`
	Limit  int32  `form:"limit" binding:"omitempty,gte=1,lte=500"`
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
