package controllers

import (
	"net/http"
	"somdeep-demo-app/src/user/interfaces"
	"somdeep-demo-app/src/user/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService interfaces.UserService
}

func NewUserController(userService interfaces.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (s *UserController) GetUsersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		response, err := s.userService.GetUsers(recordPerPage, page, startIndex)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}

func (s *UserController) GetUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		response, err := s.userService.GetUser(userId)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}

func (s *UserController) AddUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		// convert the JSON data coming from FE to something that golang understands

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Error occured while binding JSON"})
			return
		}

		response, err := s.userService.AddUser(user)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}

func (s *UserController) UpdateUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		userId := c.Param("user_id")

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Error occured while binding JSON"})
			return
		}

		response, err := s.userService.UpdateUser(userId, user)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}

func (s *UserController) DeleteUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		response, err := s.userService.DeleteUser(userId)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}
