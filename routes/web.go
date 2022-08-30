package routes

import (
	"goravel/app/http/controllers"

	"github.com/gin-gonic/gin"
	"github.com/goravel/framework/support/facades"
)

func Web() {
	facades.Route.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Hello": "Goravel",
		})
	})

	facades.Route.GET("/user", controllers.UserController{}.Show)
	facades.Route.GET("/search", controllers.SearchController{}.Search)
}
