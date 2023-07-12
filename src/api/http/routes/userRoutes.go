package routes

import (
	"somdeep-demo-app/src/api/http/controllers"
	"somdeep-demo-app/src/user/interfaces"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine, userService interfaces.UserService) {
	// Create controller instances with the userService dependency
	userController := controllers.NewUserController(userService)

	incomingRoutes.GET("/users", userController.GetUsersHandler())
	incomingRoutes.GET("/users/:user_id", userController.GetUserHandler())
	incomingRoutes.POST("/users", userController.AddUserHandler())
	incomingRoutes.PATCH("/users/:user_id", userController.UpdateUserHandler())
	incomingRoutes.DELETE("/users/:user_id", userController.DeleteUserHandler())
}
