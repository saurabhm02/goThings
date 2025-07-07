package services

import (
	"context"
	"go-auth/internals/models"
	"go-auth/internals/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	db *mongo.Collection
}

func NewUserService(db *mongo.Collection) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) RegisterUser(user models.User) (models.Response, error) {
	log.Printf("Request to register a new user! userName: %s, email: %s", user.UserName, user.Email)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if user.UserName == "" || user.Email == "" || user.Password == "" {
		log.Println("All required fields are missing")
		return models.Response{
			Status:  "Error",
			Message: "Please fill in all required fields",
		}, nil
	}

	var existingUser models.User

	query := primitive.M{"email": user.Email}
	err := s.db.FindOne(ctx, query).Decode(&existingUser)

	if err == nil {
		log.Println("User already exists with email:", existingUser.Email)
		return models.Response{
			Status:  "Conflict",
			Message: "User already exists",
		}, nil
	} else if err != mongo.ErrNoDocuments {
		log.Println("error while finding user from the db:", err)
		return models.Response{
			Status:  "Error",
			Message: "Database error while finding user",
		}, err
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("No user found with this email")
		return models.Response{
			Status:  "Error",
			Message: "No user found with this email",
		}, err
	}

	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = s.db.InsertOne(ctx, user)
	if err != nil {
		log.Fatalln("error while inserting user:", err)
		return models.Response{
			Status:  "Error",
			Message: "Error while inserting user",
		}, err
	}

	log.Println("User created successfully!")
	return models.Response{
		Status:  "Success",
		Message: "User created successfully!",
		Data:    user,
	}, nil

}

func (s *UserService) LoginUser(email, password string) (models.Response, error) {
	log.Printf("Request to login a user! email: %s", email)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := primitive.M{"email": email}
	var existingUser models.User
	err := s.db.FindOne(ctx, query).Decode(&existingUser)

	if err == mongo.ErrNoDocuments {
		log.Println("No user found with this email")
		return models.Response{
			Status:  "Error",
			Message: "No user found with this email, please create a new user",
		}, nil
	} else if err != nil {
		log.Fatalln("Error while finding user from the db:", err)
		return models.Response{
			Status:  "Error",
			Message: "Database error while finding user",
		}, err
	}

	isPasswordMatch := utils.ComparePassword(password, existingUser.Password)
	if !isPasswordMatch {
		log.Println("Password is not matching")
		return models.Response{
			Status:  "Error",
			Message: "Incorrect password",
		}, nil
	}

	token, err := utils.CreateToken(existingUser.UserName)
	if err != nil {
		log.Fatalln("Error while creating token", err)
		return models.Response{
			Status:  "Error",
			Message: "Token generation failed",
		}, err
	}

	log.Println("User logged in successfully!")
	return models.Response{
		Status:  "Success",
		Message: "User logged in successfully!",
		Data: map[string]interface{}{
			"token": token,
			"user":  existingUser,
		},
	}, nil
}

func (s *UserService) GetUserByEmail(emailId string) (models.Response, error) {
	log.Println("Request to getting user by emailId")

	query := primitive.M{"email": emailId}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := s.db.FindOne(ctx, query).Decode(&user)
	if err == nil {
		log.Println("No user found with this email")
		return models.Response{
			Status:  "Error",
			Message: "No user found with this email",
		}, err
	}

	log.Println("User found successfully!")
	return models.Response{
		Status:  "Sucess",
		Message: "User found successfully!",
		Data:    user,
	}, nil
}

func (s *UserService) IsUserAdmin(email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := primitive.M{"email": email}

	var user models.User
	err := s.db.FindOne(ctx, query).Decode(&user)
	if err != nil {
		log.Println("No user found or error in fetching user:", err)
		return false
	}

	return user.Role == models.RoleAdmin
}
