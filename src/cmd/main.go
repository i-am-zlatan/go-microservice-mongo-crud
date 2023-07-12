package main

import (
	"log"
	"os"
	"somdeep-demo-app/src/api/http/routes"
	customerMongo "somdeep-demo-app/src/customer/dal/mongo"
	customerModules "somdeep-demo-app/src/customer/modules"
	"somdeep-demo-app/src/database"
	userMongo "somdeep-demo-app/src/user/dal/mongo"
	userModules "somdeep-demo-app/src/user/modules"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	// err := godotenv.Load("/Users/somdeep/Documents/Self-Projects/somdeep-demo-app/.env")
	if err != nil {
		log.Fatal("Error loading the .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	// Initialize the MongoDB client and repository
	client := database.DBinstance()
	userRepo := userMongo.NewUserRepository(client)
	userService := userModules.NewUserService(userRepo)

	customerRepo := customerMongo.NewCustomerRepository(client)
	customerService := customerModules.NewCustomerService(customerRepo, userRepo)

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router, userService)
	routes.CustomerRoutes(router, customerService)
	router.Run(":" + port)
}
