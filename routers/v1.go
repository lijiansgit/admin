package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lijiansgit/admin/controllers"
)

const (
	URI_API = "/api"
)

func GetRouters() (router *gin.Engine) {
	router = gin.Default()
	router.GET("/ping", controllers.Ping)
	routerApi := router.Group(URI_API)
	v1 := routerApi.Group("/v1")

	// user
	user := v1.Group("/user")
	{
		user.GET("/info", controllers.User.Info)
		user.GET("/list", controllers.User.List)
		user.GET("/syncLDAP", controllers.User.LDAPUserSyncDB)
		user.POST("/login", controllers.User.Login)
		user.POST("/logout", controllers.User.Logout)
		user.POST("/modifyRoles", controllers.User.ModifyRoles)
	}

	// permission
	permission := v1.Group("/permission")
	{
		permission.GET("routes", controllers.Permission.Routes)
		permission.GET("roles", controllers.Permission.Roles)
		permission.POST("role", controllers.Permission.CreateRole)
		permission.PUT("role/:key", controllers.Permission.UpdateRole)
		permission.DELETE("role/:key", controllers.Permission.DeleteRole)
	}

	return router
}
