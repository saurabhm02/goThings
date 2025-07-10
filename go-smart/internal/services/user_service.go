package services

import (
	"errors"
	"go-smart/internal/models"
	"go-smart/internal/utils"
	"strings"
	"time"

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
	log.Infof("Request to register a user! %v", user)

	ctx, cancel := utils.GetContext()
	defer cancel()
	userName, email, mobile, password := user.Name, user.Email, user.Number, user.Password
	if userName == "" || email == "" || password == "" || mobile == 0 {
		log.Warn("required fields can not be empty, please fill all the required fields!")
		return models.Response{
			Status:  "Error",
			Message: "required fields can not be empty, please fill all the required fields!",
		}, errors.New("required fields can not be empty, please fill all the required fields!")
	}

	query := primitive.M{"email": strings.ToLower(email)}
	var existingUser models.User
	err := u.db.FindOne(ctx, query).Decode(&existingUser)
	if err == nil {
		log.Warn("User already exists")
		return models.Response{
			Status:  "Error",
			Message: "User already exists with this email!",
		}, errors.New("User already exists")
	} else if err != mongo.ErrNoDocuments {
		log.Error("DB error while checking for existing user: ", err)
		return models.Response{
			Status:  "Error",
			Message: "DB error while checking for existing user",
		}, err
	}
	hashedPassword, err := utils.HashedPassword(password)
	if err != nil {
		log.Error("Error while hashing password!")
		return models.Response{
			Status:  "Error",
			Message: "Error while hashing password!",
		}, err
	}

	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err = u.db.InsertOne(ctx, user)
	if err != nil {
		log.Error("error while inserting user:", err)
		return models.Response{
			Status:  "Error",
			Message: "Error while inserting user",
		}, err
	}
	log.Info("User create successfully!")
	otpService := NewOtpService()
	if err := otpService.GenerateAndSendOtp(email, user.SendOtpToEmail); err != nil {
		log.Error("Error while generaating otp for %s", email)
		return models.Response{
			Status:  "Error",
			Message: "Error while generaating otp",
		}, err
	}
	log.Info("User created and OTP sent for verification successfully!")
	return models.Response{
		Status:  "Success",
		Message: "User created and OTP sent for verification successfully!!",
		Data:    user,
	}, nil
}
