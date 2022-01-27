package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewTeamData represents the db fields needed to create a new team.
type NewTeamData struct {
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	ShortName string            `json:"shortName" bson:"shortName"`
	Members   map[string]string `json:"members" bson:"members"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
}

// TeamData represents the db fields associated with each team.
type TeamData struct {
	ID        string            `json:"id" bson:"_id"`
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	ShortName string            `json:"shortName" bson:"shortName"`
	Members   map[string]string `json:"members" bson:"members"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
}

// AddMemberToTeam attempts to add a new member to a team.
func AddMemberToTeam(ctx context.Context, collection *mongo.Collection, teamID, memberID, role string) (bool, error) {
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
