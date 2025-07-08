package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name           string             `bson:"name" json:"name"`
	Email          string             `bson:"email" json:"email"`
	Number         int                `bson:"number" json:"number"`
	Password       string             `bson:"password,omitempty" json:"password"`
	IsVerified     bool               `bson:"isVerified" json:"isVerified"`
	VerifiedAt     *time.Time         `bson:"verifiedAt,omitempty" json:"verifiedAt,omitempty"`
	SendOtpToEmail bool               `bson:"sendOtpToEmail" json:"sendOtpToEmail"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}
