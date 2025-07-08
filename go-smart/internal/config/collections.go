package config

import "go.mongodb.org/mongo-driver/mongo"

type Collections struct {
	UserCollection *mongo.Collection
	OtpCollection  *mongo.Collection
}
