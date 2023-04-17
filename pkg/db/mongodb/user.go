package mongodb

import (
	"context"
	"fmt"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewUser creates a new user in the db.
func (m *Mongo) NewUser(ctx context.Context, userData accounts.User) (string, error) {
	m.Initialize()
	defer m.client.Disconnect(ctx)
	err := m.client.Connect(ctx)
	if err != nil {
		m.log.Errorf("could not connect to db, %+v", err)
		return "", apierr.ErrInternalServerError
	}
	collection := m.db.Collection(CollectionUsers)
	result, err := collection.InsertOne(ctx, userData)
	if err != nil {
		m.log.Error("could not create new user")
		return "", apierr.ErrInternalServerError
	}
	userOID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		m.log.Error("error getting objectID from newly inserted user")
		return "", apierr.ErrInternalServerError
	}
	return userOID.Hex(), nil
}

// Delete attempts to delete a user from the db, returning the number of deleted users.
// TODO: remove user from all users teams.
func (m *Mongo) Delete(ctx context.Context, requestData request.DeleteUser, APIKey string) (int, apierr.Error) {
	m.Initialize()
	defer m.client.Disconnect(ctx)
	err := m.client.Connect(ctx)
	if err != nil {
		m.log.Error("could not connect to db")
		return 0, apierr.NewInternalServerError()
	}
	collection := m.db.Collection(CollectionUsers)
	res, err := m.GetByID(ctx, collection, requestData.ID)
	if err != nil {
		m.log.Errorf("could not find user to delete:  %v", err)
		return 0, apierr.NewBadRequestError("could not find user to delete")
	}
	userData, err := m.DecodeUser(res)
	if err != nil {
		m.log.Errorf("could not decode user: %v", err)
		return 0, apierr.NewBadRequestError("could not find user to delete")
	}
	ok := auth.CheckHash(userData.Password, requestData.Password)
	if !ok {
		m.log.Errorf("could not delete user - password incorrect: %v", err)
		return 0, apierr.NewWrongCredentialsError("password incorrect")
	}
	result, err := m.deleteUserFromDB(ctx, collection, requestData.ID)
	if err != nil {
		m.log.Errorf("could not delete user: %v", err)
		return 0, apierr.NewInternalServerError()
	}
	if result.DeletedCount == 0 {
		m.log.Error("no users deleted")
		return 0, apierr.NewBadRequestError("error: could not remove cmd")
	}
	return int(result.DeletedCount), nil
}

// deleteUserFromDB takes a given userID and removes the user from the database.
func (m *Mongo) deleteUserFromDB(ctx context.Context, collection *mongo.Collection, userID string) (*mongo.DeleteResult, error) {
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		m.log.Error("could not get ObjectID from hex")
		return nil, err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	result, err := collection.DeleteOne(ctx, filter, opts)
	if err != nil {
		m.log.Errorf("could not delete user: %v", err)
		return nil, err
	}
	return result, nil
}

// UserAlreadyExists checks the db for a user with given email and returns whether they already exist or not.
func (m *Mongo) UserAlreadyExists(ctx context.Context, email string) (bool, error) {
	m.Initialize()
	defer m.client.Disconnect(ctx)
	err := m.client.Connect(ctx)
	if err != nil {
		m.log.Error("couldn't connect to db")
		return false, apierr.ErrInternalServerError
	}
	collection := m.db.Collection(CollectionUsers)
	return m.DataAlreadyExists(ctx, collection, "email", email), nil
}

// GetUserByAPIKey retrieves a user from the db based on their APIKey.
func (m *Mongo) GetUserByAPIKey(ctx context.Context, APIKey string) (accounts.User, error) {
	m.Initialize()
	err := m.client.Connect(ctx)
	if err != nil {
		m.log.Error("could not connect to db")
	}
	defer m.client.Disconnect(ctx)
	collection := m.db.Collection(CollectionUsers)
	res := m.GetByKey(ctx, collection, "api_key", APIKey)
	return m.DecodeUser(res)
}

// GetUserByEmail checks the users credentials returns the user if password is correct.
func (m *Mongo) GetUserByEmail(ctx context.Context, email string) (accounts.User, error) {
	m.Initialize()
	defer m.client.Disconnect(ctx)
	err := m.client.Connect(ctx)
	if err != nil {
		m.log.Error("couldn't connect to db")
		return accounts.User{}, apierr.NewInternalServerError()
	}
	collection := m.db.Collection(CollectionUsers)
	res := m.GetByKey(ctx, collection, "email", email)
	return m.DecodeUser(res)
}

// GetAllCmds uses req info to get all users current cmds from the db.
func (m *Mongo) GetAllCmds(ctx context.Context, APIKey string) (map[string]string, apierr.Error) {
	m.Initialize()
	defer m.client.Disconnect(ctx)
	err := m.client.Connect(ctx)
	if err != nil {
		m.log.Error("couldn't connect to db")
		return nil, apierr.NewInternalServerError()
	}
	collection := m.db.Collection(CollectionUsers)
	res := m.GetByKey(ctx, collection, "api_key", APIKey)
	user, err := m.DecodeUser(res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.log.Error("couldn't find user with given APIKey")
			return nil, apierr.NewBadRequestError("could not find user")
		}
		m.log.Error("error getting user from db")
		return nil, apierr.NewInternalServerError()
	}
	return user.Cmds, nil
}

// AddCmd attempts to either add or update a cmd for the user, returning the number
// of updated cmds.
func (m *Mongo) AddCmd(ctx context.Context, requestData request.AddCmd, APIKey string) (int, apierr.Error) {
	m.Initialize()
	defer m.client.Disconnect(ctx)
	err := m.client.Connect(ctx)
	if err != nil {
		m.log.Error("couldn't connect to db")
		return 0, apierr.NewInternalServerError()
	}
	collection := m.db.Collection(CollectionUsers)
	result, err := m.addCmdToUser(ctx, collection, requestData)
	if err != nil {
		m.log.Errorf("could not add cmd to user: %v", err)
		return 0, apierr.NewInternalServerError()
	}
	var numUpdated int
	if int(result.UpsertedCount) >= int(result.ModifiedCount) {
		numUpdated = int(result.UpsertedCount)
	} else {
		numUpdated = int(result.ModifiedCount)
	}
	return numUpdated, nil
}

// addCmdToUser takes a given username along with the cmd and URL to set and adds the data to their bookmarks.
func (m *Mongo) addCmdToUser(ctx context.Context, collection *mongo.Collection, requestData request.AddCmd) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	filter, err := primitive.ObjectIDFromHex(requestData.ID)
	if err != nil {
		m.log.Error("could not get ObjectID from Hex")
		return nil, err
	}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: fmt.Sprintf("cmds.%s", requestData.Cmd), Value: requestData.URL}}}}
	result, err := collection.UpdateByID(ctx, filter, update, opts)
	if err != nil {
		m.log.Errorf("could not get update user by id: %v", err)
		return nil, err
	}
	return result, nil
}

func (m *Mongo) AddCmdByAPIKey(ctx context.Context, requestData request.AddCmd, APIKey string) (int, apierr.Error) {
	m.Initialize()
	defer m.client.Disconnect(ctx)
	err := m.client.Connect(ctx)
	if err != nil {
		m.log.Error("couldn't connect to db")
		return 0, apierr.NewInternalServerError()
	}
	collection := m.db.Collection(CollectionUsers)
	opts := UpdateEmbedOptions{
		FilterKey:   "api_key",
		FilterValue: APIKey,
		Embedded:    "cmds",
		Key:         requestData.Cmd,
		Value:       requestData.URL,
	}
	_, err = m.UpdateEmbedByField(ctx, collection, opts)
	if err != nil {
		m.log.Errorf("could not add cmd to user: %v", err)
		return 0, apierr.NewInternalServerError()
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.log.Error("couldn't find user with given APIKey")
			return 0, apierr.NewBadRequestError("could not find user")
		}
		m.log.Error("error getting user from db")
		return 0, apierr.NewInternalServerError()
	}
	return 1, nil
}

func (m *Mongo) addCmdToUserByAPIKey(ctx context.Context, collection *mongo.Collection, requestData request.AddCmd) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	filter, err := primitive.ObjectIDFromHex(requestData.ID)
	if err != nil {
		m.log.Error("could not get ObjectID from Hex")
		return nil, err
	}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: fmt.Sprintf("cmds.%s", requestData.Cmd), Value: requestData.URL}}}}
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		m.log.Errorf("could not get update user by id: %v", err)
		return nil, err
	}
	return result, nil
}

// DeleteCmd attempts to either rempve a cmd from the user, returning the number
// of updated cmds.
func (m *Mongo) DeleteCmd(ctx context.Context, requestData request.DeleteCmd, APIKey string) (int, apierr.Error) {
	m.Initialize()
	defer m.client.Disconnect(ctx)
	err := m.client.Connect(ctx)
	if err != nil {
		m.log.Error("could not connect to db")
		return 0, apierr.NewInternalServerError()
	}
	collection := m.db.Collection(CollectionUsers)
	result, err := m.removeUserCmd(ctx, collection, requestData.ID, requestData.Cmd)
	if err != nil {
		m.log.Errorf("couldn't remove cmd from user: %v", err)
		return 0, apierr.NewInternalServerError()
	}
	return int(result.ModifiedCount), nil
}

// removeUserCmd takes a given username along with the cmd and removes the cmd from their bookmarks.
func (m *Mongo) removeUserCmd(ctx context.Context, collection *mongo.Collection, userID, cmd string) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(false)
	filter, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		m.log.Error("could not get ObjectID from Hex")
		return nil, err
	}
	update := bson.D{primitive.E{Key: "$unset", Value: bson.D{primitive.E{Key: fmt.Sprintf("cmds.%s", cmd), Value: ""}}}}
	result, err := collection.UpdateByID(ctx, filter, update, opts)
	if err != nil {
		m.log.Errorf("could not get remove user cmd by ID: %v", err)
		return nil, err
	}
	return result, nil
}
