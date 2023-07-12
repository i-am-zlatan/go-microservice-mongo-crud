package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID          primitive.ObjectID `bson:"_id"`
	User_id     string             `json:"user_id"`
	Customer_id string             `json:"customer_id"`
	First_name  *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name   *string            `json:"last_name" validate:"required,min=2,max=100"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
}
