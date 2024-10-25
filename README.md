# go-jwt-gorm-gin-rsa

## Postgres
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres -d postgres:16.4-bullseye

## Routes and API Interactions.
http://localhost:9000/users
Json Payload.
{
  "username": "vineeth",
	"email":   "vs@here.at",
	"password": "xxxxxxxx",
	"Roles": ["admin"],
	"Groups": ["admin"]
}

Login:
http://localhost:9000/users/login
{
	"email":   "vs@here.at",
	"password": "xxxxxxxx",
	"forceTokenGen": false
}

Protected Routs:
http://localhost:9000/roles/

http://localhost:9000/groups/


The all routes with Gin:
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