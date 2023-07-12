package interfaces

import "somdeep-demo-app/src/customer/models"

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    any    `json:"data"`
}

type CustomerService interface {
	GetAllCustomers(recordPerPage int, page int, startIndex int) (response Response, err error)
	GetCustomersByUserId(userId string, recordPerPage int, page int, startIndex int) (response Response, err error)
	GetCustomerByCustomerId(userId string, customerId string) (response Response, err error)
	AddCustomerByUserId(userId string, customer models.Customer) (response Response, err error)
	UpdateCustomerByCustomerId(userId string, customerId string, customer models.Customer) (response Response, err error)
	DeleteCustomerByCustomerId(customerId string) (response Response, err error)
	DeleteCustomersByUserId(userId string) (response Response, err error)
}
