package services

import (
	"context"
	"errors"
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
		log.Println("No user found with this email while registering a new user")
		return models.Response{
			Status:  "Error",
			Message: "No user found with this email while registering a new user",
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
		log.Println("No user found with this email while login user")
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

	token, err := utils.CreateToken(existingUser.UserName, existingUser.Role)
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
		log.Println("No user found with this email while fetching user")
		return models.Response{
			Status:  "Error",
			Message: "No user found with this email while fetching user",
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

func (s *UserService) ChangeRoleFromAdminToUser(user models.User) (models.Response, error) {
	log.Println("Request to change the role from the admin to user")

	query := primitive.M{"email": user.Email}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := s.db.FindOne(ctx, query).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		log.Println("No user found with this email while searching for updating role")
		return models.Response{
			Status:  "Error",
			Message: "No user found with this email while searching for updating role",
		}, nil
	} else if err != nil {
		log.Println("Database error while finding user:", err)
		return models.Response{
			Status:  "Error",
			Message: "Database error while finding user",
		}, err
	}
	if existingUser.Role != models.RoleAdmin {
		log.Println("user is not an admin, sorry can't update the role", err)
		return models.Response{
			Status:  "Error",
			Message: "user is not an admin, sorry can't update the role",
		}, errors.New("user is not an admin, sorry can't update the role")
	}
	update := primitive.M{
		"$set": primitive.M{
			"role":      models.RoleUser,
			"updatedAt": time.Now(),
		},
	}
	_, err = s.db.UpdateOne(ctx, query, update)
	if err != nil {
		log.Println("Error while updating user role:", err)
		return models.Response{
			Status:  "Error",
			Message: "Error while updating user role",
		}, err
	}
	log.Println("User role changed to user successfully!")
	return models.Response{
		Status:  "Success",
		Message: "User role changed to user successfully!",
	}, nil
}

func (s *UserService) DeleteUser(emailId string) (models.Response, error) {
	log.Println("Request to delete user with email:", emailId)

	query := primitive.M{"email": emailId}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, check if user exists
	var existingUser models.User
	err := s.db.FindOne(ctx, query).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		log.Println("No user found with this email while searching for deleting user")
		return models.Response{
			Status:  "Error",
			Message: "No user found with this email while searching for deleting user",
		}, nil
	} else if err != nil {
		log.Println("Database error while finding user:", err)
		return models.Response{
			Status:  "Error",
			Message: "Database error while finding user",
		}, err
	}

	result, err := s.db.DeleteOne(ctx, query)
	if err != nil {
		log.Println("Error while deleting user:", err)
		return models.Response{
			Status:  "Error",
			Message: "Failed to delete user",
		}, err
	}

	if result.DeletedCount == 0 {
		log.Println("No user deleted - already removed?")
		return models.Response{
			Status:  "Error",
			Message: "No user was deleted",
		}, nil
	}

	log.Println("User deleted successfully!")
	return models.Response{
		Status:  "Success",
		Message: "User deleted successfully",
	}, nil
}
