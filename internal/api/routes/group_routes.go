package routes

import (
	"web_socket/internal/api/handlers"
	"web_socket/internal/common/middleware"

	"github.com/gin-gonic/gin"
)

type GroupRoutes struct {
	handler *handlers.GroupHandler
}

func NewGroupRoutes(handler *handlers.GroupHandler) *GroupRoutes {
	return &GroupRoutes{
		handler: handler,
	}
}
func (groupRoute *GroupRoutes) Register(r *gin.RouterGroup) {
	group := r.Group("/groups")
	group.Use(middleware.AuthMiddleware())
	{
		group.POST("/create", groupRoute.handler.CreateGroup)
		group.GET("/getAll", groupRoute.handler.GetAllGroups)
		group.PUT("/:uuid", groupRoute.handler.UpdateGroup)
		group.DELETE("/:uuid", groupRoute.handler.SoftDeleteGroup)
		group.POST("/:uuid/leave", groupRoute.handler.LeaveGroup)
		group.GET("/:uuid/members", groupRoute.handler.GetGroupMembers)
		group.GET("/:uuid/members/:user_uuid", groupRoute.handler.GetMemberInfo)
		group.POST("/:uuid/members", groupRoute.handler.AddMemberToGroup)
		group.PUT("/:uuid/members/", groupRoute.handler.UpdateMemberRole)
		group.DELETE("/:uuid/members/:user_uuid", groupRoute.handler.RemoveMember)
	}
}
