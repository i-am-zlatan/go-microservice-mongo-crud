package mongo

import (
	"context"
	"errors"
	"somdeep-demo-app/src/database"
	"somdeep-demo-app/src/user/interfaces"
	"somdeep-demo-app/src/user/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) interfaces.UserRepository {
	userCollection := database.OpenCollection(client, "user")
	return &userRepository{
		userCollection: userCollection,
	}
}

// var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func (r *userRepository) GetAllUsers(startIndex int, recordPerPage int, ctx context.Context) (result *mongo.Cursor, err error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}}, {Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}}, {Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}
	result, err = r.userCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage})

	return result, err
}

func (r *userRepository) GetUserByUserId(ctx context.Context, userId string) (user models.User, result error) {
	result = r.userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	return user, result
}

func (r *userRepository) CountDocumentBasedOnKey(ctx context.Context, user models.User, key string) (count int64, err error) {
	filter := bson.M{}
	switch key {
	case "email":
		filter["email"] = user.Email
	case "phone":
		filter["phone"] = user.Phone
	default:
		return 0, errors.New("unsupported key")
	}

	count, err = r.userCollection.CountDocuments(ctx, filter)
	return count, err
}

func (r *userRepository) AddUserToMongoDb(ctx context.Context, user models.User) (insertErr error) {
	_, insertErr = r.userCollection.InsertOne(ctx, user)
	return insertErr
}

func (r *userRepository) UpdateOneUserByUserId(ctx context.Context, opt options.UpdateOptions, filter primitive.M, updateObject primitive.D) (result *mongo.UpdateResult, err error) {
	result, err = r.userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObject},
		},
		&opt,
	)

	return result, err
}

func (r *userRepository) DeleteOneUserByUserId(ctx context.Context, filter primitive.M) (result *mongo.DeleteResult, err error) {
	result, err = r.userCollection.DeleteOne(ctx, filter)
	return result, err
}
