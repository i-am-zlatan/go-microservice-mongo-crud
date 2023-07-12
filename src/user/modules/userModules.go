package modules

import (
	"context"
	"log"
	"net/http"
	"somdeep-demo-app/src/user/interfaces"
	"somdeep-demo-app/src/user/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type userService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(userRepository interfaces.UserRepository) interfaces.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetUsers(recordPerPage int, page int, startIndex int) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var res interfaces.Response

	result, err := s.userRepository.GetAllUsers(startIndex, recordPerPage, ctx)
	defer cancel()
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "error occured while listing user items"
		res.Data = nil
		return res, err
	}
	var allusers []bson.M
	if err = result.All(ctx, &allusers); err != nil {
		log.Fatal(err)
	}

	if len(allusers) == 0 {
		res.Status = http.StatusInternalServerError
		res.Error = "NA"
		res.Message = "No Records Found"
		res.Data = nil
		return res, err
	}

	// c.JSON(http.StatusOK, allusers)
	res.Status = http.StatusOK
	res.Error = "NA"
	res.Message = "Records Fetched Successfully"
	res.Data = allusers
	return res, err
}

func (s *userService) GetUser(userId string) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res interfaces.Response
	var user models.User
	// err = userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	user, err = s.userRepository.GetUserByUserId(ctx, userId)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"message": "Error occured while fetching documents", "error": err.Error()})
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "Error occured while listing user items"
		res.Data = nil
		return res, err
	}
	// c.JSON(http.StatusOK, user)
	res.Status = http.StatusOK
	res.Error = "NA"
	res.Message = "Record Fetched Successfully"
	res.Data = user
	return res, err
}

func (s *userService) AddUser(user models.User) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res interfaces.Response

	// validate the data based on user struct

	validationError := validate.Struct(user)

	if validationError != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
		res.Status = http.StatusBadRequest
		res.Error = err.Error()
		res.Message = "Validation Error"
		res.Data = nil
		return res, err
	}

	// we will check whether the "email" has already been used by another user or not

	count, err := s.userRepository.CountDocumentBasedOnKey(ctx, user, "email")

	if err != nil {
		log.Panic(err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the mail"})
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "Error occured while checking for e-mail"
		res.Data = nil
		return res, err
	}

	// hash the password - HashPassword()

	password := HashPassword(*user.Password)
	user.Password = &password

	// we will check whether the "phone_number" has already been used by another user or not

	count, err = s.userRepository.CountDocumentBasedOnKey(ctx, user, "phone")

	if err != nil {
		log.Panic(err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "Error occured while checking for phone number"
		res.Data = nil
		return res, err
	}

	if count > 0 {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
		res.Status = http.StatusInternalServerError
		res.Error = "NA"
		res.Message = "User with this e-mail or phone already exists"
		res.Data = nil
		return res, err
	}

	// create some extra details for the user object - basically fillers (created_at, updated_at and ID)

	user.Created_at = time.Now()
	user.Updated_at = time.Now()
	user.ID = primitive.NewObjectID()
	user.User_id = uuid.New().String()

	insertErr := s.userRepository.AddUserToMongoDb(ctx, user)
	if insertErr != nil {
		// msg := "User item was not created"
		// c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		res.Status = http.StatusInternalServerError
		res.Error = insertErr.Error()
		res.Message = "User item was not created"
		res.Data = nil
		return res, err
	}
	// c.JSON(http.StatusOK, gin.H{"message": "User added successfully", "userID": user.User_id})
	res.Status = http.StatusOK
	res.Error = "NA"
	res.Message = "User Added Successfully"
	res.Data = "user_id: " + user.User_id
	return res, err
}

func (s *userService) UpdateUser(userId string, user models.User) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res interfaces.Response
	var updateObject primitive.D

	if user.First_name != nil {
		updateObject = append(updateObject, bson.E{Key: "first_name", Value: user.First_name})
	}

	if user.Last_name != nil {
		updateObject = append(updateObject, bson.E{Key: "last_name", Value: user.Last_name})
	}

	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObject = append(updateObject, bson.E{Key: "updated_at", Value: user.Updated_at})

	upsert := false

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	filter := bson.M{"user_id": userId}

	result, err := s.userRepository.UpdateOneUserByUserId(ctx, opt, filter, updateObject)

	if result.ModifiedCount == 0 {
		// c.JSON(http.StatusNotFound, gin.H{"message": "User not found or is already deleted"})
		res.Status = http.StatusNotFound
		res.Error = "NA"
		res.Message = "User not found or is already deleted"
		res.Data = nil
		return res, err
	}

	if err != nil {
		// msg := "User update failed"
		// c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "User update failed"
		res.Data = nil
		return res, err
	}
	// c.JSON(http.StatusOK, result)

	res.Status = http.StatusOK
	res.Error = "NA"
	res.Message = "User updated successfully"
	res.Data = result
	return res, err
}

func (s *userService) DeleteUser(userId string) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res interfaces.Response
	filter := bson.M{"user_id": userId}

	result, err := s.userRepository.DeleteOneUserByUserId(ctx, filter)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "Failed to delete user"
		res.Data = nil
		return res, nil
	}

	if result.DeletedCount == 0 {
		// c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		res.Status = http.StatusNotFound
		res.Error = "NA"
		res.Message = "User not found or is already deleted"
		res.Data = nil
		return res, nil
	}

	// c.JSON(http.StatusOK, gin.H{"message": "User deleted", "userId": userId})
	res.Status = http.StatusOK
	res.Error = "NA"
	res.Message = "User deleted successfully"
	res.Data = "user_id: " + userId
	return res, nil
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "login or password is incorrect"
		check = false
	}

	return check, msg
}
