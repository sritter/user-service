package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "user-service/docs" // swag docs
	"user-service/user"
)

// @title User Service API
// @version 1.0
// @description REST API for user CRUD operations
// @host localhost:8080
// @BasePath /
func Run() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes := r.Group("/users")
	{
		routes.POST("", user.CreateUser)
		routes.GET("", user.ListUsers)
		routes.GET(":id", user.GetUser)
		routes.PUT(":id", user.UpdateUser)
		routes.DELETE(":id", user.DeleteUser)
	}

	r.Run(":8080")
}
