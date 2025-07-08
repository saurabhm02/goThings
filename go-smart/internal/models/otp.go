package models

import (
	"time"
)

type OTP struct {
	UserId    string    `bson:"userId" json:"userId"`
	Otp       string    `bson:"otp" json:"otp"`
	Type      string    `bson:"type" json:"type"`
	ExpireAt  time.Time `bson:"expireAt" json:"expireAt"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
