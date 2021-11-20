package models

import (
	"context"
	"math/rand"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignUpReq struct {
	Name     *string `json:"name"`
	Password *string `json:"password"`
}

func UserFieldAlreadyExists(ctx context.Context, collection *mongo.Collection, key, value string) bool {
	var result bson.M
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: key, Value: value}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
	}
	return true
}

func GenerateAPIKey() string {
	chars := strings.Split("QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm0123456789-", "")
	rand.Shuffle(len(chars), func(a, b int) {
		chars[a], chars[b] = chars[b], chars[a]
	})
	length := rand.Intn(len(chars)-10) + 10
	key := strings.Join(chars[:length], "")
	return key
}
