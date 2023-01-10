package mongodb

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type tokenDocument struct {
	ID     string `bson:"_id"`
	APIKey string `bson:"api_key"`
	Token  string `bson:"token"`
}

func (m *Mongo) GetRefreshTokenByAPIKey(ctx context.Context, APIKey string) (string, error) {
	m.Initialize()
	err := m.client.Connect(ctx)
	if err != nil {
		m.log.Error("could not connect to db")
		return "", apierr.ErrInternalServerError
	}
	defer m.client.Disconnect(ctx)
	collection := m.db.Collection(CollectionTokens)
	res := m.GetByKey(ctx, collection, "api_key", APIKey)
	var token tokenDocument
	err = res.Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", apierr.ErrNotFound
		}
		return "", apierr.ErrInternalServerError
	}
	return token.Token, nil
}

func (m *Mongo) NewRefreshToken(ctx context.Context, APIKey, refreshToken string) error {
	m.Initialize()
	err := m.client.Connect(ctx)
	if err != nil {
		m.log.Error("could not connect to db")
		return apierr.ErrInternalServerError
	}
	defer m.client.Disconnect(ctx)
	collection := m.db.Collection(CollectionTokens)
	update := bson.M{"$set": bson.M{"token": refreshToken}}
	options := options.Update().SetUpsert(true)
	res, err := collection.UpdateOne(ctx, bson.M{"api_key": APIKey}, update, options)
	if err != nil || (res.ModifiedCount+res.UpsertedCount < 1) {
		m.log.Error("could not update db with refresh token")
		return apierr.ErrInternalServerError
	}
	return nil
}
