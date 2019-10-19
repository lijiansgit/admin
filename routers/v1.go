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
	user := v1.Group("/user")
	{
		user.GET("/info", controllers.User.Info)
		user.POST("/login", controllers.User.Login)
		user.POST("/logout", controllers.User.Logout)
	}

	return router
}
