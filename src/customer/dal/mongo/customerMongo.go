package mongo

import (
	"context"
	"somdeep-demo-app/src/customer/interfaces"
	"somdeep-demo-app/src/customer/models"
	"somdeep-demo-app/src/database"

	// userMongo "somdeep-demo-app/src/user/dal/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type customerRepository struct {
	customerCollection *mongo.Collection
}

func NewCustomerRepository(client *mongo.Client) interfaces.CustomerRepository {
	customerCollection := database.OpenCollection(client, "customer")
	return &customerRepository{
		customerCollection: customerCollection,
	}
}

func (r *customerRepository) GetAllCustomers(startIndex int, recordPerPage int, ctx context.Context) (result *mongo.Cursor, err error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}}, {Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}}, {Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}
	result, err = r.customerCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage})

	return result, err
}

func (r *customerRepository) GetCustomersByUserId(userId string, startIndex int, recordPerPage int, ctx context.Context) (result *mongo.Cursor, err error) {
	matchStage := bson.D{{Key: "$match", Value: bson.M{"user_id": userId}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}}, {Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}}, {Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}
	result, err = r.customerCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage})

	return result, err
}

func (r *customerRepository) GetCustomerByCustomerId(ctx context.Context, customerId string) (customer models.Customer, result error) {
	result = r.customerCollection.FindOne(ctx, bson.M{"customer_id": customerId}).Decode(&customer)
	return customer, result
}

func (r *customerRepository) AddCustomerToMongoDb(ctx context.Context, customer models.Customer) (insertErr error) {
	_, insertErr = r.customerCollection.InsertOne(ctx, customer)
	return insertErr
}

func (r *customerRepository) UpdateCustomerByCustomerId(ctx context.Context, opt options.UpdateOptions, filter primitive.M, updateObject primitive.D) (result *mongo.UpdateResult, err error) {
	result, err = r.customerCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObject},
		},
		&opt,
	)
	return result, err
}

func (r *customerRepository) DeleteCustomerByCustomerId(customerId string, ctx context.Context, filter primitive.M) (result *mongo.DeleteResult, err error) {
	result, err = r.customerCollection.DeleteOne(ctx, filter)
	return result, err
}

func (r *customerRepository) DeleteCustomersByUserId(customerId string, ctx context.Context, filter primitive.M) (result *mongo.DeleteResult, err error) {
	result, err = r.customerCollection.DeleteMany(ctx, filter)
	return result, err
}
