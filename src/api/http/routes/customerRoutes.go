package routes

import (
	"somdeep-demo-app/src/api/http/controllers"
	"somdeep-demo-app/src/customer/interfaces"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(incomingRoutes *gin.Engine, customerService interfaces.CustomerService) {
	// Create controller instances with the userService dependency
	customerController := controllers.NewCustomerController(customerService)
	incomingRoutes.GET("/customers", customerController.GetCustomersHandler())
	incomingRoutes.GET("/users/:user_id/customers", customerController.GetCustomersByUserIdHandler())
	incomingRoutes.GET("/users/:user_id/customers/:customer_id", customerController.GetCustomerByCustomerIdHandler())
	incomingRoutes.POST("/users/:user_id/customers", customerController.AddCustomerByUserIdHandler())
	incomingRoutes.PATCH("/users/:user_id/customers/:customer_id", customerController.UpdateCustomerByCustomerIdHandler())
	incomingRoutes.DELETE("/users/:user_id/customers/:customer_id", customerController.DeleteCustomerByCustomerIdHandler())
	incomingRoutes.DELETE("/users/:user_id/customers", customerController.DeleteCustomersByUserId())
}
