package main

import (
	"jwt/controller"
	"jwt/initializers"
	"jwt/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.InitialierEnvVariable()
	initializers.InitiazeDB()
	initializers.MigrateDB()
	initializers.SeedRoles()
}

func main() {
	r := gin.Default()

	userGroup := r.Group("/users")
	{
		userGroup.POST("/", controller.CreateUser)
		userGroup.POST("/login", controller.LoginUser)
		userGroup.GET("/:id", controller.GetUser)
		userGroup.PUT("/:id", controller.UpdateUser)
		userGroup.DELETE("/:id", controller.DeleteUser)
		userGroup.GET("/", controller.ListUsers)
	}

	groupGroup := r.Group("/groups")
	groupGroup.Use(middleware.AdminRequired())
	{
		groupGroup.POST("/", controller.CreateGroup)
		groupGroup.GET("/:id", controller.GetGroup)
		groupGroup.PUT("/:id", controller.UpdateGroup)
		groupGroup.DELETE("/:id", controller.DeleteGroup)
		groupGroup.GET("/", controller.ListGroups)
	}

	roleGroup := r.Group("/roles")
	roleGroup.Use(middleware.AdminRequired())
	{
		roleGroup.POST("/", controller.CreateRole)
		roleGroup.GET("/:id", controller.GetRole)
		roleGroup.PUT("/:id", controller.UpdateRole)
		roleGroup.DELETE("/:id", controller.DeleteRole)
		roleGroup.GET("/", controller.ListRoles)
	}

	r.Run()
}
