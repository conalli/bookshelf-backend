package user

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/conalli/bookshelf-backend/auth/password"
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// New checks whether a username alreadys exists in the db. If not, a new user
// is created based upon the request data.
func New(reqCtx context.Context, requestData requests.CredentialsRequest) (string, string, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	userExists := DataAlreadyExists(ctx, &collection, "name", requestData.Name)

	if userExists {
		return "", "", errors.NewBadRequestError(fmt.Sprintf("error creating new user; user with name %v already exists", requestData.Name))
	}
	apiKey := GenerateAPIKey()
	for DataAlreadyExists(ctx, &collection, "apiKey", apiKey) {
		apiKey = GenerateAPIKey()
	}
	hashedPassword, err := password.HashPassword(requestData.Password)
	if err != nil {
		log.Println("error hashing password")
		return "", "", errors.NewInternalServerError()
	}
	newUserData := NewUserData{
		Name:      requestData.Name,
		Password:  hashedPassword,
		APIKey:    apiKey,
		Bookmarks: map[string]string{},
		Teams:     map[string]string{},
	}
	res, err := collection.InsertOne(ctx, newUserData)
	if err != nil {
		log.Printf("error creating new user with data: \n username: %v\n password: %v", requestData.Name, requestData.Password)
		return "", "", errors.NewInternalServerError()
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("error getting objectID from newly inserted user")
		return "", "", errors.NewInternalServerError()
	}
	return oid.Hex(), apiKey, nil
}

// DataAlreadyExists attempts to find a user based on a given key-value pair, returning wether they
// already exist in the db or not.
func DataAlreadyExists(ctx context.Context, collection *mongo.Collection, key, value string) bool {
	var result bson.M
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: key, Value: value}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
	}
	return true
}

// GenerateAPIKey generates a random URL-safe string of random length for use as an API key.
// TODO: Refactor method for generating keys.
func GenerateAPIKey() string {
	chars := strings.Split("QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm0123456789-", "")
	rand.Shuffle(len(chars), func(a, b int) {
		chars[a], chars[b] = chars[b], chars[a]
	})
	length := rand.Intn(len(chars)-10) + 10
	key := strings.Join(chars[:length], "")
	return key
}
