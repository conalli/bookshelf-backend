package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/conalli/bookshelf-backend/pkg/accounts"
	"github.com/conalli/bookshelf-backend/pkg/db"
	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/password"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewTeam checks whether a team name alreadys exists in the db. If not, a new team
// is created based upon the request data.
func (m *Mongo) NewTeam(ctx context.Context, requestData accounts.NewTeamRequest) (string, errors.ApiErr) {
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
	hashedPassword, err := password.HashPassword(requestData.TeamPassword)
	if err != nil {
		log.Printf("couldnt hash password: %+v\n", err)
		return "", errors.NewInternalServerError()
	}
	newTeamData := accounts.Team{
		Name:      requestData.Name,
		Password:  hashedPassword,
		ShortName: requestData.ShortName,
		Members:   map[string]string{requestData.ID: accounts.RoleAdmin},
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
			Value:       accounts.RoleAdmin,
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

// DeleteTeam removes a team from the db and updates each users account to reflect it.
// TODO: improve validation
func (m *Mongo) DeleteTeam(ctx context.Context, requestData accounts.DelTeamRequest) (int, errors.ApiErr) {
	reqCtx, cancelFunc := db.ReqContextWithTimeout(ctx)
	defer cancelFunc()
	m.Initialize()
	err := m.client.Connect(reqCtx)
	if err != nil {
		log.Printf("couldn't connect to db on new team, %+v\n", err)
		return 0, errors.NewInternalServerError()
	}
	defer m.client.Disconnect(reqCtx)
	res, err := m.SessionWithTransaction(reqCtx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		teamCollection := m.db.Collection(CollectionTeams)
		res, err := GetByID(sessCtx, teamCollection, requestData.TeamID)
		if err != nil {
			log.Printf("Couldnt find team with id: %s to delete.\n", requestData.TeamID)
			return nil, errors.NewBadRequestError("error couldn't find team with ID: " + requestData.TeamID)
		}
		teamData, err := DecodeTeam(res)
		if err != nil {
			log.Println("Couldn't decode team from result.")
			return nil, errors.NewInternalServerError()
		}
		if !password.CheckHashedPassword(teamData.Password, requestData.TeamPassword) {
			log.Println("error: wrong password when attempting to delete team")
			return nil, errors.NewWrongCredentialsError("incorrect team password")
		}
		result, err := deleteTeamFromDB(sessCtx, teamCollection, requestData.TeamID)
		if err != nil {
			log.Printf("error deleting team %s from db\n", requestData.TeamID)
			return nil, err
		}
		userCollection := m.db.Collection(CollectionUsers)
		memberIDs := getMemberIDsFromTeam(teamData.Members)
		update, err := deleteTeamFromUsers(sessCtx, userCollection, teamData.ID, memberIDs)
		if err != nil {
			log.Printf("couldn't update team members to remove team %s\n", requestData.TeamID)
			return nil, err
		}
		if update.MatchedCount == 0 {
			log.Printf("Update couldn't match any members for team: %s\n", requestData.TeamID)
			return nil, errors.NewInternalServerError()
		}
		if update.ModifiedCount == 0 {
			log.Printf("Couldn't update any members of team: %s\n", requestData.TeamID)
			return nil, errors.NewInternalServerError()
		}
		return int(result.DeletedCount), nil
	})
	if err != nil {
		log.Printf("error trying to delete team -> %v", err)
		return 0, errors.NewInternalServerError()
	}
	v, ok := res.(int)
	if !ok {
		log.Printf("result of transaction type %T, wanted int\n", res)
		return 0, errors.NewInternalServerError()
	}
	return v, nil
}

func deleteTeamFromDB(ctx context.Context, collection *mongo.Collection, teamID string) (*mongo.DeleteResult, error) {
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	id, err := primitive.ObjectIDFromHex(teamID)
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

func getMemberIDsFromTeam(userIDs map[string]string) []string {
	ids := make([]string, len(userIDs))
	i := 0
	for k := range userIDs {
		ids[i] = k
		i++
	}
	return ids
}

func deleteTeamFromUsers(ctx context.Context, collection *mongo.Collection, teamID string, userIDs []string) (*mongo.UpdateResult, error) {
	var oids []primitive.ObjectID
	for _, userID := range userIDs {
		oid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Println("couldnt get objectID from hex")
			return nil, err
		}
		oids = append(oids, oid)
	}
	filter := bson.D{primitive.E{Key: "_id", Value: bson.D{primitive.E{Key: "$in", Value: oids}}}}
	update := bson.D{primitive.E{Key: "$unset", Value: bson.D{primitive.E{Key: fmt.Sprintf("teams.%s", teamID), Value: ""}}}}
	res, err := collection.UpdateMany(ctx, filter, update)
	return res, err
}

// AddMember checks whether a username alreadys exists in the db. If not, a new user
// is created based upon the request data.
func (m *Mongo) AddMember(ctx context.Context, requestData accounts.AddMemberRequest) (bool, errors.ApiErr) {
	reqCtx, cancelFunc := db.ReqContextWithTimeout(ctx)
	defer cancelFunc()
	m.Initialize()
	err := m.client.Connect(reqCtx)
	if err != nil {
		log.Printf("couldn't connect to db on new team, %+v", err)
	}
	defer m.client.Disconnect(reqCtx)
	res, err := m.SessionWithTransaction(reqCtx, func(sessCtx mongo.SessionContext) (interface{}, error) {
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

// TODO: add validation for role
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

// DeleteSelf takes a request and attemps to remove a member from the given team.
func (m *Mongo) DeleteSelf(ctx context.Context, requestData accounts.DelSelfRequest) (bool, errors.ApiErr) {
	reqCtx, cancelFunc := db.ReqContextWithTimeout(ctx)
	defer cancelFunc()
	m.Initialize()
	err := m.client.Connect(reqCtx)
	if err != nil {
		log.Printf("couldn't connect to db on new team, %+v", err)
	}
	defer m.client.Disconnect(reqCtx)
	res, err := m.SessionWithTransaction(reqCtx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		userCollection := m.db.Collection(CollectionUsers)
		update := UpdateEmbedOptions{
			FilterKey:   "_id",
			FilterValue: requestData.ID,
			Embedded:    "teams",
			Key:         requestData.TeamID,
			Value:       "",
		}
		res, err := UpdateEmbedByField(sessCtx, userCollection, update)
		if err != nil {
			log.Println("couldnt remove user with id: " + requestData.ID + "from team: " + requestData.TeamID)
			return false, errors.NewBadRequestError("couldnt remove user from team")
		}
		user, err := DecodeUser(res)
		if err != nil {
			log.Printf("couldnt decode user from single update result: %+v", err)
			return false, err
		}
		teamCollection := m.db.Collection(CollectionTeams)
		ok, err := removeMemberFromTeam(sessCtx, teamCollection, requestData.TeamID, user.ID)
		if err != nil {
			log.Printf("error removing user with name: %s from team with id: %s\n error: %+v\n", requestData.ID, requestData.TeamID, err)
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

// DeleteMember takes a request and attemps to remove a member from the given team.
func (m *Mongo) DeleteMember(ctx context.Context, requestData accounts.DelMemberRequest) (bool, errors.ApiErr) {
	reqCtx, cancelFunc := db.ReqContextWithTimeout(ctx)
	defer cancelFunc()
	m.Initialize()
	err := m.client.Connect(reqCtx)
	if err != nil {
		log.Printf("couldn't connect to db on new team, %+v", err)
	}
	defer m.client.Disconnect(reqCtx)
	res, err := m.SessionWithTransaction(reqCtx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		userCollection := m.db.Collection(CollectionUsers)
		update := UpdateEmbedOptions{
			FilterKey:   "name",
			FilterValue: requestData.MemberName,
			Embedded:    "teams",
			Key:         requestData.TeamID,
			Value:       "",
		}
		res, err := UpdateEmbedByField(sessCtx, userCollection, update)
		if err != nil {
			log.Println("couldnt remove user with id: " + requestData.ID + "from team: " + requestData.TeamID)
			return false, errors.NewBadRequestError("couldnt remove user from team")
		}
		user, err := DecodeUser(res)
		if err != nil {
			log.Printf("couldnt decode user from single update result: %+v", err)
			return false, err
		}
		teamCollection := m.db.Collection(CollectionTeams)
		ok, err := removeMemberFromTeam(sessCtx, teamCollection, requestData.TeamID, user.ID)
		if err != nil {
			log.Printf("error removing member with name: %s from team with id: %s\n error: %+v\n", requestData.MemberName, requestData.TeamID, err)
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

func removeMemberFromTeam(ctx context.Context, collection *mongo.Collection, teamID, memberID string) (bool, error) {
	opts := options.Update().SetUpsert(true)
	filter, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		log.Printf("error getting objectid from hex: %+v\n", err)
		return false, err
	}
	update := bson.D{primitive.E{Key: "$unset", Value: bson.D{primitive.E{Key: fmt.Sprintf("members.%s", memberID), Value: ""}}}}
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

// AddTeamCmd takes request data and attempts to add a new cmd to the teams bookmarks.
// TODO: improve validation.
func (m *Mongo) AddTeamCmd(ctx context.Context, requestData accounts.AddTeamCmdRequest) (int, errors.ApiErr) {
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

// DelTeamCmd attempts to either rempve a cmd from the user, returning the number
// of updated cmds.
// TODO: improve validation
func (m *Mongo) DelTeamCmd(ctx context.Context, requestData accounts.DelTeamCmdRequest, APIKey string) (int, errors.ApiErr) {
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
