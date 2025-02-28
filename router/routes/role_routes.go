package routes

import (
	"github.com/gin-gonic/gin"
	"hrkGo/app/controllers/c_system"
	"hrkGo/app/middleware"
)

func RegisterRoleRoutes(r *gin.RouterGroup) {

	roleController := c_system.RoleController{}

	role := r.Group("/role")
	{
		authRouter := role.Use(middleware.JWTAuth())
		{
			authRouter.GET("/list", roleController.GetRoleList)
			authRouter.GET("/:roleId", roleController.SelectRoleDataById)
			authRouter.GET("/deptTree/:roleId", roleController.SelectRoleDataById)

		}
	}

}
