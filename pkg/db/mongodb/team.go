package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/conalli/bookshelf-backend/pkg/db"
	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/password"
	"github.com/conalli/bookshelf-backend/pkg/team"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Team represents the db fields associated with each team.
type Team struct {
	ID        string            `json:"id" bson:"_id,omitempty"`
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	ShortName string            `json:"shortName" bson:"shortName"`
	Members   map[string]string `json:"members" bson:"members"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
}

// New checks whether a team name alreadys exists in the db. If not, a new team
// is created based upon the request data.
func (m *Mongo) New(ctx context.Context, requestData team.NewTeamRequest) (string, errors.ApiErr) {
	reqCtx, cancelFunc := db.ReqContextWithTimeout(ctx)
	defer cancelFunc()
	m.Initialize()
	err := m.client.Connect(reqCtx)
	if err != nil {
		log.Printf("couldn't connect to db on new team, %+v", err)
	}
	defer m.client.Disconnect(reqCtx)
	teamCollection := m.db.Collection(CollectionTeams)
	teamExists := DataAlreadyExists(reqCtx, teamCollection, "name", requestData.Name)
	if teamExists {
		log.Println("team already exists")
		return "", errors.NewBadRequestError(fmt.Sprintf("error creating new user; user with name %v already exists", requestData.Name))
	}
	hashedPassword, err := password.HashPassword(requestData.Password)
	if err != nil {
		log.Printf("couldnt hash password: %+v\n", err)
		return "", errors.NewInternalServerError()
	}
	newTeamData := Team{
		Name:      requestData.Name,
		Password:  hashedPassword,
		ShortName: requestData.ShortName,
		Members:   map[string]string{requestData.ID: "admin"},
		Bookmarks: map[string]string{},
	}
	res, err := m.SessionWithTransaction(reqCtx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		result, err := teamCollection.InsertOne(sessCtx, newTeamData)
		if err != nil {
			log.Printf("couldnt insert team %+v\n", err)
			return "", errors.NewInternalServerError()
		}
		teamOID, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			log.Println("couldnt convert InsertedID to ObjectID")
			return "", errors.NewInternalServerError()
		}
		userCollection := m.db.Collection(CollectionUsers)
		opts := UpdateEmbedOptions{
			FilterKey:   "_id",
			FilterValue: requestData.ID,
			Embedded:    "teams",
			Key:         teamOID.Hex(),
			Value:       "admin",
		}
		res, err := UpdateEmbedByField(sessCtx, userCollection, opts)
		if err != nil {
			log.Printf("couldnt update embedded field: %+v\n", err)
			return nil, err
		}
		_, err = DecodeUser(res)
		if err != nil {
			log.Printf("couldnt decode user from single update result: %+v", err)
			return nil, err
		}
		return teamOID.Hex(), nil
	})
	if err != nil {
		log.Println("error during session transaction")
		return "", errors.NewInternalServerError()
	}
	v, ok := res.(string)
	if !ok {
		return "", errors.NewInternalServerError()
	}
	return v, nil
}

// AddMember checks whether a username alreadys exists in the db. If not, a new user
// is created based upon the request data.
func (m *Mongo) AddMember(ctx context.Context, requestData team.AddMemberRequest) (bool, errors.ApiErr) {
	reqCtx, cancelFunc := db.ReqContextWithTimeout(ctx)
	defer cancelFunc()
	m.Initialize()
	err := m.client.Connect(reqCtx)
	if err != nil {
		log.Printf("couldn't connect to db on new team, %+v", err)
	}
	defer m.client.Disconnect(reqCtx)
	res, err := m.SessionWithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		userCollection := m.db.Collection(CollectionUsers)
		update := UpdateEmbedOptions{
			FilterKey:   "name",
			FilterValue: requestData.MemberName,
			Embedded:    "teams",
			Key:         requestData.TeamID,
			Value:       requestData.Role,
		}
		res, err := UpdateEmbedByField(sessCtx, userCollection, update)
		if err != nil {
			log.Println("couldnt update user with name: " + requestData.MemberName)
			return false, errors.NewBadRequestError("couldnt update user with name: " + requestData.MemberName)
		}
		user, err := DecodeUser(res)
		if err != nil {
			log.Printf("couldnt decode user from single update result: %+v", err)
			return false, err
		}
		teamCollection := m.db.Collection(CollectionTeams)
		ok, err := addMemberToTeam(sessCtx, teamCollection, requestData.TeamID, user.ID, requestData.Role)
		if err != nil {
			log.Printf("error adding member with name: %s to team with id: %s\n error: %+v\n", requestData.MemberName, requestData.TeamID, err)
			return false, errors.NewInternalServerError()
		}
		return ok, nil
	})
	if err != nil {
		log.Println("error could not start db transaction")
		return false, errors.NewInternalServerError()
	}
	if v, ok := res.(bool); ok {
		if !v {
			return false, nil
		}
	} else {
		return false, nil
	}
	return true, nil
}

func addMemberToTeam(ctx context.Context, collection *mongo.Collection, teamID, memberID, role string) (bool, error) {
	opts := options.Update().SetUpsert(true)
	filter, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		log.Printf("error getting objectid from hex: %+v\n", err)
		return false, err
	}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: fmt.Sprintf("members.%s", memberID), Value: role}}}}
	result, err := collection.UpdateByID(ctx, filter, update, opts)
	if err != nil {
		log.Printf("error attempting to add user: %s to team: %s: %+v\n", memberID, teamID, err)
		return false, err
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("error attempting to add user: %s to team: %s, team was not modified\n", memberID, teamID)
		return false, nil
	}
	return true, nil
}

// AddCmdToTeam takes request data and attempts to add a new cmd to the teams bookmarks.
// TODO: improve validation.
func (m *Mongo) AddCmdToTeam(ctx context.Context, requestData team.AddTeamCmdRequest) (int, errors.ApiErr) {
	reqCtx, cancelFunc := db.ReqContextWithTimeout(ctx)
	defer cancelFunc()
	m.Initialize()
	err := m.client.Connect(reqCtx)
	if err != nil {
		log.Printf("couldn't connect to db on new team, %+v", err)
	}
	defer m.client.Disconnect(reqCtx)
	collection := m.db.Collection(CollectionTeams)
	result, err := addCmdToTeam(reqCtx, collection, requestData.ID, requestData.Cmd, requestData.URL)
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

func addCmdToTeam(ctx context.Context, collection *mongo.Collection, userID, key, value string) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	filter, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: fmt.Sprintf("bookmarks.%s", key), Value: value}}}}
	result, err := collection.UpdateByID(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DelCmdFromTeam attempts to either rempve a cmd from the user, returning the number
// of updated cmds.
// TODO: improve validation
func (m *Mongo) DelCmdFromTeam(ctx context.Context, requestData team.DelTeamCmdRequest, APIKey string) (int, errors.ApiErr) {
	reqCtx, cancelFunc := db.ReqContextWithTimeout(ctx)
	defer cancelFunc()
	m.Initialize()
	err := m.client.Connect(reqCtx)
	if err != nil {
		log.Printf("couldn't connect to db on new team, %+v", err)
	}
	defer m.client.Disconnect(reqCtx)

	teamCollection := m.db.Collection(CollectionTeams)
	result, err := removeCmdFromTeam(reqCtx, teamCollection, requestData.ID, requestData.Cmd)
	if err != nil {
		return 0, errors.NewInternalServerError()
	}
	return int(result.ModifiedCount), nil
}

func removeCmdFromTeam(ctx context.Context, collection *mongo.Collection, teamID, cmd string) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	filter, err := primitive.ObjectIDFromHex(teamID)
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
