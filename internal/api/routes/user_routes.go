package routes

import (
	"web_socket/internal/api/handlers"
	"web_socket/internal/common/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *handlers.UserHandler
}

func NewUserRoutes(handler *handlers.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}
func (userRoute *UserRoutes) Register(r *gin.RouterGroup) {
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("/:uuid", middleware.RequireSelfOrAdmin("uuid"), userRoute.handler.FindUserByUUID)
		users.DELETE("/:uuid", middleware.RequireSelfOrAdmin("uuid"), userRoute.handler.SoftDeleteUser)
		adminMod := users.Group("")
		adminMod.Use(middleware.RequirePermission(middleware.AdminRoleID))
		{
			adminMod.PUT("/:uuid/restore", userRoute.handler.RestoreUser)
			adminMod.DELETE("/:uuid/clean", userRoute.handler.HardDeleteUser)
		}
	}
}
