package services

import (
	"errors"
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

func (o *OtpService) GenerateAndSendOtp(email string, sendOtpToEmail bool) error {
	log.Infof("Generating OTP for email: %s", email)

	ctx, cancel := utils.GetContext()
	defer cancel()

	var user models.User
	err := o.userDb.FindOne(ctx, primitive.M{"email": email}).Decode(&user)
	if err != nil {
		log.Error("User not found: ", err)
		return err
	}

	if sendOtpToEmail {
		log.Infof("Sending OTP to email: %s", user.Email)
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
			return err
		}

		err = utils.SendMail(user.Email, "Your OTP Code", "Your OTP is: "+otpValue)
		if err != nil {
			log.Error("Failed to send OTP via email: ", err)
			return err
		}
	} else {
		if user.Number == 0 {
			log.Error("User phone number is missing")
			return errors.New("user phone number is missing")
		}
		log.Infof("Sending OTP to phone: ****%d", user.Number%10000)
		_, err = utils.SendSms(strconv.Itoa(user.Number))
		if err != nil {
			log.Error("Failed to send OTP via SMS: ", err)
			return err
		}
	}
	log.Info("Otp generate and send successfully!")

	return nil
}

func (o *OtpService) VerifyOtp(email string, otp string, isEmailOtp bool) (bool, error) {
	ctx, cancel := utils.GetContext()
	defer cancel()

	var user models.User
	err := o.userDb.FindOne(ctx, primitive.M{"email": email}).Decode(&user)
	if err != nil {
		log.Error("User not found: ", err)
		return false, err
	}

	if isEmailOtp {
		var storedOtp models.OTP
		err := o.otpDb.FindOne(ctx, primitive.M{
			"email": user.Email,
			"otp":   otp,
		}).Decode(&storedOtp)
		if err != nil {
			log.Error("Invalid OTP or not found: ", err)
			return false, err
		}

		if time.Now().After(storedOtp.ExpireAt) {
			log.Warn("OTP has expired")
			return false, nil
		}

		update := primitive.M{
			"$set": primitive.M{
				"isVerified": true,
				"updatedAt":  time.Now(),
			},
		}
		_, err = o.userDb.UpdateOne(ctx, primitive.M{"email": email}, update)
		if err != nil {
			log.Error("Failed to update user verification status: ", err)
			return false, err
		}

		_, _ = o.otpDb.DeleteOne(ctx, primitive.M{"email": storedOtp.Email})
		log.Info("OTP verified successfully (email)")
		return true, nil

	} else {
		if user.Number == 0 {
			log.Error("User phone number is missing")
			return false, errors.New("user phone number is missing")
		}
		err := utils.CheckOtp(strconv.Itoa(user.Number), otp)
		if err != nil {
			log.Error("Twilio OTP verification failed: ", err)
			return false, err
		}

		update := primitive.M{
			"$set": primitive.M{
				"isVerified": true,
				"updatedAt":  time.Now(),
			},
		}
		_, err = o.userDb.UpdateOne(ctx, primitive.M{"email": email}, update)
		if err != nil {
			log.Error("Failed to update user verification status: ", err)
			return false, err
		}

		log.Info("OTP verified successfully (SMS)")
		return true, nil
	}
}
