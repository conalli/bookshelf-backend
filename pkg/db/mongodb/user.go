package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/password"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewUser is a func.
func (m *Mongo) NewUser(ctx context.Context, requestData accounts.SignUpRequest) (accounts.User, errors.ApiErr) {
	m.Initialize()
	err := m.client.Connect(ctx)
	if err != nil {
		log.Printf("couldn't connect to db on new user, %+v", err)
	}
	defer m.client.Disconnect(ctx)
	collection := m.db.Collection(CollectionUsers)
	userExists := DataAlreadyExists(ctx, collection, "name", requestData.Name)
	if userExists {
		log.Println("user already exists")
		return accounts.User{}, errors.NewBadRequestError(fmt.Sprintf("error creating new user; user with name %v already exists", requestData.Name))
	}
	APIKey, err := accounts.GenerateAPIKey()
	if err != nil {
		log.Println("error generating uuid")
		return accounts.User{}, errors.NewInternalServerError()
	}
	hashedPassword, err := password.HashPassword(requestData.Password)
	if err != nil {
		log.Println("error hashing password")
		return accounts.User{}, errors.NewInternalServerError()
	}
	signUpData := accounts.User{
		Name:      requestData.Name,
		Password:  hashedPassword,
		APIKey:    APIKey,
		Bookmarks: map[string]string{},
		Teams:     map[string]string{},
	}
	res, err := collection.InsertOne(ctx, signUpData)
	if err != nil {
		log.Printf("error creating new user with data: \n username: %v\n password: %v", requestData.Name, requestData.Password)
		return accounts.User{}, errors.NewInternalServerError()
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("error getting objectID from newly inserted user")
		return accounts.User{}, errors.NewInternalServerError()
	}
	newUserData := accounts.User{
		ID:     oid.Hex(),
		Name:   requestData.Name,
		APIKey: APIKey,
	}
	return newUserData, nil
}

// GetUserByName checks the users credentials returns the user if password is correct.
func (m *Mongo) GetUserByName(ctx context.Context, requestData accounts.LogInRequest) (accounts.User, error) {
	m.Initialize()
	err := m.client.Connect(ctx)
	if err != nil {
		log.Println("couldn't connect to db on login")
	}
	defer m.client.Disconnect(ctx)
	collection := m.db.Collection(CollectionUsers)
	res := GetByKey(ctx, collection, "name", requestData.Name)
	return DecodeUser(res)
}

// GetTeams uses user id to get all users teams from the db.
func (m *Mongo) GetTeams(ctx context.Context, APIKey string) ([]accounts.Team, errors.ApiErr) {
	m.Initialize()
	err := m.client.Connect(ctx)
	if err != nil {
		log.Println("couldn't connect to db on login")
	}
	defer m.client.Disconnect(ctx)
	res, err := m.SessionWithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		userCollection := m.db.Collection(CollectionUsers)
		res := GetByKey(sessCtx, userCollection, "APIKey", APIKey)
		currUser, err := DecodeUser(res)
		if err != nil {
			log.Printf("error getting user by APIKey: %s -> %+v\n", APIKey, err)
			return nil, err
		}
		teamIDs, err := convertIDs(currUser.Teams)
		if err != nil {
			log.Printf("error converting teams to ids: error -> %+v\n", err)
			return nil, err
		}
		teamCollection := m.db.Collection(CollectionTeams)
		filter := bson.M{"_id": bson.M{"$in": teamIDs}}
		opts := options.Find()
		teamCursor, err := teamCollection.Find(sessCtx, filter, opts)
		if err != nil {
			log.Printf("error converting teams to ids -> %+v\n", err)
			return nil, err
		}
		defer teamCursor.Close(sessCtx)
		var teams []accounts.Team
		for teamCursor.Next(sessCtx) {
			var currTeam accounts.Team
			if err := teamCursor.Decode(&currTeam); err != nil {
				log.Printf("error could not get team from found teams -> %+v\n", err)
				return nil, err
			}
			teams = append(teams, currTeam)
		}
		return teams, nil
	})
	if err != nil {
		log.Printf("error could not get data from transaction -> %+v\n", err)
		return nil, errors.NewInternalServerError()
	}
	teams, ok := res.([]accounts.Team)
	if !ok {
		log.Println("error could not assert type []Team")
		return nil, errors.NewInternalServerError()
	}
	return teams, nil
}

func convertIDs(teams map[string]string) ([]primitive.ObjectID, error) {
	output := make([]primitive.ObjectID, len(teams))
	for id := range teams {
		res, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		output = append(output, res)
	}
	return output, nil
}

// GetAllCmds uses req info to get all users current cmds from the db.
func (m *Mongo) GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.ApiErr) {
	m.Initialize()
	err := m.client.Connect(ctx)
	if err != nil {
		log.Println("couldn't connect to db on login")
	}
	defer m.client.Disconnect(ctx)
	collection := m.db.Collection(CollectionUsers)

	res := GetByKey(ctx, collection, "APIKey", APIKey)
	user, err := DecodeUser(res)
	if err != nil {
		return nil, errors.ParseGetUserError(APIKey, err)
	}
	return user.Bookmarks, nil
}

// AddCmd attempts to either add or update a cmd for the user, returning the number
// of updated cmds.
func (m *Mongo) AddCmd(ctx context.Context, requestData accounts.AddCmdRequest, APIKey string) (int, errors.ApiErr) {
	m.Initialize()
	err := m.client.Connect(ctx)
	if err != nil {
		log.Println("couldn't connect to db on login")
	}
	defer m.client.Disconnect(ctx)
	collection := m.db.Collection(CollectionUsers)

	result, err := addCmdToUser(ctx, collection, requestData)
	if err != nil {
		return 0, errors.NewInternalServerError()
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
func addCmdToUser(ctx context.Context, collection *mongo.Collection, requestData accounts.AddCmdRequest) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	filter, err := primitive.ObjectIDFromHex(requestData.ID)
	if err != nil {
		return nil, err
	}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: fmt.Sprintf("bookmarks.%s", requestData.Cmd), Value: requestData.URL}}}}
	result, err := collection.UpdateByID(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteCmd attempts to either rempve a cmd from the user, returning the number
// of updated cmds.
func (m *Mongo) DeleteCmd(ctx context.Context, requestData accounts.DelCmdRequest, APIKey string) (int, errors.ApiErr) {
	m.Initialize()
	err := m.client.Connect(ctx)
	if err != nil {
		log.Println("couldn't connect to db on login")
	}
	defer m.client.Disconnect(ctx)
	collection := m.db.Collection(CollectionUsers)
	result, err := removeUserCmd(ctx, collection, requestData.ID, requestData.Cmd)
	if err != nil {
		return 0, errors.NewInternalServerError()
	}
	return int(result.ModifiedCount), nil
}

// removeUserCmd takes a given username along with the cmd and removes the cmd from their bookmarks.
func removeUserCmd(ctx context.Context, collection *mongo.Collection, userID, cmd string) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	filter, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	update := bson.D{primitive.E{Key: "$unset", Value: bson.D{primitive.E{Key: fmt.Sprintf("bookmarks.%s", cmd), Value: ""}}}}
	result, err := collection.UpdateByID(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete attempts to delete a user from the db, returning the number of deleted users.
// TODO: remove user from all users teams.
func (m *Mongo) Delete(ctx context.Context, requestData accounts.DelUserRequest, APIKey string) (int, errors.ApiErr) {
	m.Initialize()
	err := m.client.Connect(ctx)
	if err != nil {
		log.Println("couldn't connect to db on login")
	}
	defer m.client.Disconnect(ctx)
	collection := m.db.Collection(CollectionUsers)
	res, err := GetByID(ctx, collection, requestData.ID)
	if err != nil {
		log.Printf("error deleting user: couldn't find user -> %v", err)
		return 0, errors.NewBadRequestError("could not find user to delete")
	}
	userData, err := DecodeUser(res)
	if err != nil {
		log.Printf("error decoding user -> %v", err)
		return 0, errors.NewBadRequestError("could not find user to delete")
	}
	ok := password.CheckHashedPassword(userData.Password, requestData.Password)
	if !ok {
		log.Printf("error deleting user: password incorrect -> %v", err)
		return 0, errors.NewWrongCredentialsError("password incorrect")
	}
	result, err := deleteUserFromDB(ctx, collection, requestData.ID)
	if err != nil {
		log.Printf("error deleting user: error -> %v", err)
		return 0, errors.NewInternalServerError()
	}
	if result.DeletedCount == 0 {
		log.Printf("could not remove user... maybe user:%s doesn't exists?", requestData.Name)
		return 0, errors.NewBadRequestError("error: could not remove cmd")
	}
	return int(result.DeletedCount), nil
}

// deleteUserFromDB takes a given userID and removes the user from the database.
func deleteUserFromDB(ctx context.Context, collection *mongo.Collection, userID string) (*mongo.DeleteResult, error) {
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	result, err := collection.DeleteOne(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetUserByAPIKey retrieves a user from the db based on their APIKey.
func (m *Mongo) GetUserByAPIKey(ctx context.Context, APIKey string) (accounts.User, error) {
	m.Initialize()
	err := m.client.Connect(ctx)
	if err != nil {
		log.Println("couldjnt connect to db on search")
	}
	defer m.client.Disconnect(ctx)
	collection := m.db.Collection(CollectionUsers)
	res := GetByKey(ctx, collection, "APIKey", APIKey)
	return DecodeUser(res)
}
