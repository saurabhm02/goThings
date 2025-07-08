package models

import (
	"time"
)

type OTP struct {
	Email     string    `bson:"email" json:"email"`
	Number    int       `bson:"number,omitempty" json:"number,omitempty"`
	Otp       string    `bson:"otp" json:"otp"`
	ExpireAt  time.Time `bson:"expireAt" json:"expireAt"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
