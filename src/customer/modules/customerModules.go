package modules

import (
	"context"
	"log"
	"net/http"
	"time"

	"somdeep-demo-app/src/customer/interfaces"
	"somdeep-demo-app/src/customer/models"
	userInterfaces "somdeep-demo-app/src/user/interfaces"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate = validator.New()

type customerService struct {
	customerRepository interfaces.CustomerRepository
	userRepository     userInterfaces.UserRepository
}

func NewCustomerService(customerRepository interfaces.CustomerRepository, userRepository userInterfaces.UserRepository) interfaces.CustomerService {
	return &customerService{
		customerRepository: customerRepository,
		userRepository:     userRepository,
	}
}

func (s *customerService) GetAllCustomers(recordPerPage int, page int, startIndex int) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var res interfaces.Response

	result, err := s.customerRepository.GetAllCustomers(startIndex, recordPerPage, ctx)
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

func (s *customerService) GetCustomersByUserId(userId string, recordPerPage int, page int, startIndex int) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res interfaces.Response

	_, err = s.userRepository.GetUserByUserId(ctx, userId)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "The user associated with customer is not present or is deleted"
		res.Data = nil
		return res, err
	}
	result, err := s.customerRepository.GetCustomersByUserId(userId, startIndex, recordPerPage, ctx)
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

func (s *customerService) GetCustomerByCustomerId(userId string, customerId string) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res interfaces.Response
	var customer models.Customer
	// err = userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	_, err = s.userRepository.GetUserByUserId(ctx, userId)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "The user associated with customer is not present or is deleted"
		res.Data = nil
		return res, err
	}
	customer, err = s.customerRepository.GetCustomerByCustomerId(ctx, customerId)
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
	res.Data = customer
	return res, err
}

func (s *customerService) AddCustomerByUserId(userId string, customer models.Customer) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res interfaces.Response

	_, err = s.userRepository.GetUserByUserId(ctx, userId)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "The user associated with customer is not present or is deleted"
		res.Data = nil
		return res, err
	}

	// validate the data based on user struct

	validationError := validate.Struct(customer)

	if validationError != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
		res.Status = http.StatusBadRequest
		res.Error = err.Error()
		res.Message = "Validation Error"
		res.Data = nil
		return res, err
	}

	customer.Created_at = time.Now()
	customer.Updated_at = time.Now()
	customer.ID = primitive.NewObjectID()
	customer.Customer_id = uuid.New().String()
	customer.User_id = userId

	insertErr := s.customerRepository.AddCustomerToMongoDb(ctx, customer)
	if insertErr != nil {
		// msg := "User item was not created"
		// c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		res.Status = http.StatusInternalServerError
		res.Error = insertErr.Error()
		res.Message = "Customer item was not created"
		res.Data = nil
		return res, err
	}
	// c.JSON(http.StatusOK, gin.H{"message": "User added successfully", "userID": user.User_id})
	res.Status = http.StatusOK
	res.Error = "NA"
	res.Message = "Customer Added Successfully"
	res.Data = "user_id: " + customer.User_id + " & customer_id: " + customer.Customer_id
	return res, err
}

func (s *customerService) UpdateCustomerByCustomerId(userId string, customerId string, customer models.Customer) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res interfaces.Response
	var updateObject primitive.D

	_, err = s.userRepository.GetUserByUserId(ctx, userId)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "The user associated with customer is not present or is deleted"
		res.Data = nil
		return res, err
	}

	if customer.First_name != nil {
		updateObject = append(updateObject, bson.E{Key: "first_name", Value: customer.First_name})
	}

	if customer.Last_name != nil {
		updateObject = append(updateObject, bson.E{Key: "last_name", Value: customer.Last_name})
	}

	customer.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObject = append(updateObject, bson.E{Key: "updated_at", Value: customer.Updated_at})

	upsert := false

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	filter := bson.M{"customer_id": customerId}

	result, err := s.customerRepository.UpdateCustomerByCustomerId(ctx, opt, filter, updateObject)

	if result.ModifiedCount == 0 {
		// c.JSON(http.StatusNotFound, gin.H{"message": "User not found or is already deleted"})
		res.Status = http.StatusNotFound
		res.Error = "NA"
		res.Message = "Customer not found or is already deleted"
		res.Data = nil
		return res, err
	}

	if err != nil {
		// msg := "User update failed"
		// c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "Customer update failed"
		res.Data = nil
		return res, err
	}
	// c.JSON(http.StatusOK, result)

	res.Status = http.StatusOK
	res.Error = "NA"
	res.Message = "Customer updated successfully"
	res.Data = result
	return res, err
}

func (s *customerService) DeleteCustomerByCustomerId(customerId string) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res interfaces.Response
	filter := bson.M{"customer_id": customerId}

	result, err := s.customerRepository.DeleteCustomerByCustomerId(customerId, ctx, filter)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "Failed to delete customer"
		res.Data = nil
		return res, nil
	}

	if result.DeletedCount == 0 {
		// c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		res.Status = http.StatusNotFound
		res.Error = "NA"
		res.Message = "Customer not found or is already deleted"
		res.Data = nil
		return res, nil
	}

	// c.JSON(http.StatusOK, gin.H{"message": "User deleted", "userId": userId})
	res.Status = http.StatusOK
	res.Error = "NA"
	res.Message = "Customer deleted successfully"
	res.Data = "customer_id: " + customerId
	return res, nil
}

func (s *customerService) DeleteCustomersByUserId(userId string) (response interfaces.Response, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res interfaces.Response
	filter := bson.M{"user_id": userId}

	result, err := s.customerRepository.DeleteCustomersByUserId(userId, ctx, filter)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		res.Status = http.StatusInternalServerError
		res.Error = err.Error()
		res.Message = "Failed to delete customer"
		res.Data = nil
		return res, nil
	}

	if result.DeletedCount == 0 {
		// c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		res.Status = http.StatusNotFound
		res.Error = "NA"
		res.Message = "Customer not found or is already deleted"
		res.Data = nil
		return res, nil
	}

	// c.JSON(http.StatusOK, gin.H{"message": "User deleted", "userId": userId})
	res.Status = http.StatusOK
	res.Error = "NA"
	res.Message = "Customer deleted successfully"
	res.Data = "user_id: " + userId
	return res, nil
}

// func (s *customerService) CheckUserExistsOrNot(ctx context.Context, userId string) (res interfaces.Response, err error) {
// 	_, err = s.userRepository.GetUserByUserId(ctx, userId)
// 	if err != nil {
// 		res.Status = http.StatusInternalServerError
// 		res.Error = err.Error()
// 		res.Message = "The user associated with customer is not present or is deleted"
// 		res.Data = nil
// 		return res, err
// 	}
// 	return res, err
// }
