package services

import (
	"errors"
	"go-smart/internal/models"
	"go-smart/internal/utils"

	"github.com/gofiber/fiber/v2/log"
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

func (u *UserService) RegisterUser(user models.User) (models.Response, error) {
	log.Info("Request to register a user! %v", user)

	ctx, canel := utils.GetContext()
	defer canel()
	userName, email, mobile, password := user.Name, user.Email, user.Number, user.Password
	if userName == "" || email == "" || mobile == "" || password == "" {
		log.Warn("required fields can not be empty, please fill all the required fields!")
		return models.Response{
			Status:  "Error",
			Message: "required fields can not be empty, please fill all the required fields!",
		}, errors.New("required fields can not be empty, please fill all the required fields!")
	}

	var existingUser *models.User
	query := primitive.M{"email": email}
	err := u.db.FindOne(ctx, query).Decode(&existingUser)
	if err != nil {
		log.Warn("user is already exist!, use different email to create a new user!")
		return models.Response{
			Status:  "Error",
			Message: "user is already exist!, use different email to create a new user!",
		}, errors.New("user is already exist!, use different email to create a new user!")
	}

}
