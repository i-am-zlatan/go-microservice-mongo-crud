package interfaces

import (
	"context"
	"somdeep-demo-app/src/user/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	GetAllUsers(startIndex int, recordPerPage int, ctx context.Context) (*mongo.Cursor, error)
	GetUserByUserId(ctx context.Context, userId string) (models.User, error)
	CountDocumentBasedOnKey(ctx context.Context, user models.User, key string) (int64, error)
	AddUserToMongoDb(ctx context.Context, user models.User) error
	UpdateOneUserByUserId(ctx context.Context, opt options.UpdateOptions, filter primitive.M, updateObject primitive.D) (*mongo.UpdateResult, error)
	DeleteOneUserByUserId(ctx context.Context, filter primitive.M) (*mongo.DeleteResult, error)
}
