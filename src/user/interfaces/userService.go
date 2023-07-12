package interfaces

import "somdeep-demo-app/src/user/models"

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    any    `json:"data"`
}

type UserService interface {
	GetUsers(recordPerPage int, page int, startIndex int) (response Response, err error)
	GetUser(userId string) (response Response, err error)
	AddUser(user models.User) (response Response, err error)
	UpdateUser(userId string, user models.User) (response Response, err error)
	DeleteUser(userId string) (response Response, err error)
}
