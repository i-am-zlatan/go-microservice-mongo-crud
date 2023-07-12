package interfaces

import (
	"context"
	"somdeep-demo-app/src/customer/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Res struct {
	Result any `json:"result"`
}

type CustomerRepository interface {
	GetAllCustomers(startIndex int, recordPerPage int, ctx context.Context) (result *mongo.Cursor, err error)
	GetCustomersByUserId(userId string, startIndex int, recordPerPage int, ctx context.Context) (result *mongo.Cursor, err error)
	GetCustomerByCustomerId(ctx context.Context, customerId string) (customer models.Customer, result error)
	AddCustomerToMongoDb(ctx context.Context, customer models.Customer) (insertErr error)
	UpdateCustomerByCustomerId(ctx context.Context, opt options.UpdateOptions, filter primitive.M, updateObject primitive.D) (result *mongo.UpdateResult, err error)
	DeleteCustomerByCustomerId(customerId string, ctx context.Context, filter primitive.M) (result *mongo.DeleteResult, err error)
	DeleteCustomersByUserId(customerId string, ctx context.Context, filter primitive.M) (result *mongo.DeleteResult, err error)
}
