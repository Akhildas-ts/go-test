package repository

import (
	"context"
	"fmt"
	"lock/database"
	"lock/domain"
	"lock/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckingEmailValidation(email string) (*domain.User, error) {

	var user domain.User

	// Get the users collection
	collection := database.DB.Collection("users")

	// Search for the user with the given email
	filter := bson.M{"email": email}

	// Execute the query
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("error finding user: %v", err)
	}

	fmt.Println("user details:", user)
	return &user, nil

}

func ChechingPhoneExist(phone string) (*domain.User, error) {
	var user domain.User

	// Get the users collection
	collection := database.DB.Collection("users")

	// Define the filter to search for the user with the given phone
	filter := bson.M{"phone": phone}

	// Execute the query to find the user
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		// If no user is found, return nil
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		// If an error occurs during the query, return the error
		return nil, fmt.Errorf("error finding user: %v", err)
	}

	// Return the user
	return &user, nil
}

func SignupInsert(user models.SignupDetail) (models.SignupDetailResponse, error) {

	var signupRes models.SignupDetailResponse

	// Get the users collection
	collection := database.DB.Collection("users")

	// Define the document to insert
	doc := bson.M{
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"email":     user.Email,
		"phone":     user.Phone,
		"password":  user.Password,
	}

	// Insert the document
	result, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		return models.SignupDetailResponse{}, fmt.Errorf("error inserting user: %v", err)
	}

	// Set the response
	signupRes.ID = result.InsertedID.(int)
	signupRes.FirstName = user.FirstName
	signupRes.LastName = user.LastName
	signupRes.Email = user.Email
	signupRes.Phone = user.Phone

	fmt.Println("signup inserted data's are :", signupRes)
	return signupRes, nil
}

func FindUserDetailByEmail(user models.LoginDetails) (models.UserLoginResponse, error) {

	var userDetails models.UserLoginResponse

	// Get the users collection
	collection := database.DB.Collection("users")

	// Define the filter to search for the user with the given email
	filter := bson.M{"email": user.Email, "blocked": false}

	// Execute the query to find the user
	err := collection.FindOne(context.Background(), filter).Decode(&userDetails)
	if err != nil {
		// If no user is found, return an empty response
		if err == mongo.ErrNoDocuments {
			return models.UserLoginResponse{}, nil
		}
		// If an error occurs during the query, return the error
		return models.UserLoginResponse{}, fmt.Errorf("error finding user: %v", err)
	}

	// Return the user details
	return userDetails, nil
}
