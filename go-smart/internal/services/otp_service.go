package services

import (
	"go-smart/internal/config"
	"go-smart/internal/models"
	"go-smart/internal/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OtpService struct {
	otpDb  *mongo.Collection
	userDb *mongo.Collection
}

func NewOtpService() *OtpService {
	return &OtpService{
		otpDb:  config.GetDB().Collection("OTP"),
		userDb: config.GetDB().Collection("User"),
	}
}

func (o *OtpService) GenerateAndSendOtp(email string, sendOtpToEmail bool) (string, error) {
	log.Infof("Generating OTP for email: %s", email)

	ctx, cancel := utils.GetContext()
	defer cancel()

	var user models.User
	err := o.userDb.FindOne(ctx, primitive.M{"email": email}).Decode(&user)
	if err != nil {
		log.Error("User not found: ", err)
		return "", err
	}

	if sendOtpToEmail {
		otpValue := strconv.Itoa(utils.GenerateOtp())

		otp := models.OTP{
			Email:     user.Email,
			Number:    user.Number,
			Otp:       otpValue,
			ExpireAt:  time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = o.otpDb.InsertOne(ctx, otp)
		if err != nil {
			log.Error("Failed to store OTP: ", err)
			return "", err
		}

		err = utils.SendMail(user.Email, "Your OTP Code", "Your OTP is: "+otpValue)
		if err != nil {
			log.Error("Failed to send OTP via email: ", err)
			return "", err
		}
	} else {
		_, err = utils.SendSms(strconv.Itoa(user.Number))
		if err != nil {
			log.Error("Failed to send OTP via SMS: ", err)
			return "", err
		}
	}

	return "OTP sent successfully", nil
}
