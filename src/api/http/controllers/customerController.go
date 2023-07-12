package controllers

import (
	"net/http"
	"somdeep-demo-app/src/customer/interfaces"
	"somdeep-demo-app/src/customer/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	customerService interfaces.CustomerService
}

func NewCustomerController(customerService interfaces.CustomerService) *CustomerController {
	return &CustomerController{
		customerService: customerService,
	}
}

func (s *CustomerController) GetCustomersHandler() gin.HandlerFunc {
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
		response, err := s.customerService.GetAllCustomers(recordPerPage, page, startIndex)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}

func (s *CustomerController) GetCustomersByUserIdHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId := c.Param("user_id")
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		response, err := s.customerService.GetCustomersByUserId(userId, recordPerPage, page, startIndex)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}

func (s *CustomerController) GetCustomerByCustomerIdHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerId := c.Param("customer_id")
		userId := c.Param("user_id")

		response, err := s.customerService.GetCustomerByCustomerId(userId, customerId)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}

func (s *CustomerController) AddCustomerByUserIdHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId := c.Param("user_id")
		var customer models.Customer

		// convert the JSON data coming from FE to something that golang understands

		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Error occured while binding JSON"})
			return
		}

		response, err := s.customerService.AddCustomerByUserId(userId, customer)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}

func (s *CustomerController) UpdateCustomerByCustomerIdHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer models.Customer

		userId := c.Param("user_id")
		customerId := c.Param("customer_id")

		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Error occured while binding JSON"})
			return
		}

		response, err := s.customerService.UpdateCustomerByCustomerId(userId, customerId, customer)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}

func (s *CustomerController) DeleteCustomerByCustomerIdHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerId := c.Param("customer_id")

		response, err := s.customerService.DeleteCustomerByCustomerId(customerId)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}

func (s *CustomerController) DeleteCustomersByUserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		response, err := s.customerService.DeleteCustomersByUserId(userId)

		if err != nil {
			c.JSON(response.Status, response)
			return
		}

		// Return the users as the response
		c.JSON(response.Status, response)
	}
}
