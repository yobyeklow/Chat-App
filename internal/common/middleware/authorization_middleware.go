package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

const AdminRoleID int32 = 2

func RequirePermission(allowRoles ...int32) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, authorized := isRoleAuthorized(ctx, allowRoles)
		if !authorized {
			return
		}
		ctx.Next()
	}
}
func RequireSelfOrAdmin(uuidParamName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if isAdmin(ctx) {
			ctx.Next()
			return
		}
		userUUID, exists := ctx.Get("user_uuid")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User UUID not found"})
			return
		}
		targetUUID := ctx.Param(uuidParamName)
		if targetUUID == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "UUID parameter missing"})
			return
		}
		if userUUID.(string) != targetUUID {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You cannot perform this action!"})
			return
		}
		ctx.Next()
	}
}
func isAdmin(ctx *gin.Context) bool {
	userRole, exists := ctx.Get("user_role")
	if !exists {
		return false
	}
	role, ok := userRole.(int32)
	if !ok {
		return false
	}
	return role == AdminRoleID
}
func isRoleAuthorized(ctx *gin.Context, allowedRoles []int32) (int32, bool) {
	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return 0, false // Added missing return to prevent panic
	}
	role := userRole.(int32)
	if slices.Contains(allowedRoles, role) {
		return role, true
	}
	ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
	return role, false
}
