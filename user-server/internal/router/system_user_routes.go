package router

import (
	"hivemtk-user/internal/controller"
	"hivemtk-user/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupSystemUserRoutes(auth *gin.RouterGroup) {
	ctrl := controller.NewSystemUserAdminController()
	// 读（成员查看）：admin/tenant_admin/team_admin 均可（数据范围 TenantUserScoped 过滤）
	memberRead := auth.Group("/system/users", middleware.TeamAdminAuthMiddleware())
	{
		memberRead.GET("", ctrl.GetUsers)
		memberRead.GET("/:id", ctrl.GetByID)
	}
	// 写（成员管理）：admin + tenant_admin
	memberWrite := auth.Group("/system/users", middleware.TenantAdminAuthMiddleware())
	{
		memberWrite.POST("", ctrl.Create)
		memberWrite.PUT("/:id", ctrl.Update)
		memberWrite.DELETE("/:id", ctrl.Delete)
	}
}

func setupRoleRoutes(auth *gin.RouterGroup) {
	ctrl := controller.NewRoleController()
	admin := auth.Group("/system/roles", middleware.TenantAdminAuthMiddleware())
	{
		admin.GET("", ctrl.ListRoles)
		admin.GET("/:code", ctrl.GetRole)
		admin.GET("/:code/members", ctrl.ListMembers)
	}
}

func setupPermissionRoutes(auth *gin.RouterGroup) {
	ctrl := controller.NewPermissionController()
	admin := auth.Group("/system/permissions", middleware.TenantAdminAuthMiddleware())
	{
		admin.GET("/audit-logs", ctrl.ListAuditLogs)
		admin.PUT("/:id/enabled", ctrl.SetEnabled)
		admin.PUT("/:id/password", ctrl.ResetPassword)
	}
}
