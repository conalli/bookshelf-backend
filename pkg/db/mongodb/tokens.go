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
	collection := m.db.Collection(CollectionTokens)
	res := m.GetByKey(ctx, collection, "api_key", APIKey)
	var token tokenDocument
	err := res.Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", apierr.ErrNotFound
		}
		return "", apierr.ErrInternalServerError
	}
	return token.Token, nil
}

func (m *Mongo) NewRefreshToken(ctx context.Context, APIKey, refreshToken string) error {
	collection := m.db.Collection(CollectionTokens)
	update := bson.M{"$set": bson.M{"token": refreshToken}}
	options := options.Update().SetUpsert(true)
	res, err := collection.UpdateOne(ctx, bson.M{"api_key": APIKey}, update, options)
	if err != nil || (res.ModifiedCount+res.UpsertedCount < 1) {
		m.log.Errorf("could not update db with refresh token: %+v", err)
		return apierr.ErrInternalServerError
	}
	return nil
}

func (m *Mongo) DeleteRefreshToken(ctx context.Context, APIKey string) (int64, error) {
	collection := m.db.Collection(CollectionTokens)
	res, err := collection.DeleteOne(ctx, bson.M{"api_key": APIKey})
	if err != nil {
		m.log.Errorf("could not remove refesh token from db: %+v", err)
		return 0, nil
	}
	return res.DeletedCount, nil
}
