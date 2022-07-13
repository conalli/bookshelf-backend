package dbtest

// import (
// 	"context"

// 	"github.com/conalli/bookshelf-backend/pkg/accounts"
// 	"github.com/conalli/bookshelf-backend/pkg/errors"
// )

// func (t *dbtest) NewTeam(ctx context.Context, requestData request.NewTeam) (string, errors.APIErr) {
// 	found := t.dataAlreadyExists(requestData.Name, "teams")
// 	if found {
// 		return "", errors.NewBadRequestError("error creating new user; user with name " + requestData.Name + " already exists")
// 	}
// 	team := accounts.Team{
// 		ID:        "999",
// 		Name:      requestData.Name,
// 		Password:  requestData.TeamPassword,
// 		ShortName: requestData.ShortName,
// 		Members:   map[string]string{requestData.ID: accounts.RoleAdmin},
// 		Bookmarks: map[string]string{},
// 	}
// 	t.Teams[team.ID] = team
// 	return team.ID, nil
// }

// func (t *dbtest) DeleteTeam(ctx context.Context, requestData accounts.request.DeleteTeam) (int, errors.APIErr) {
// 	val, found := t.Teams[requestData.TeamID]
// 	if !found {
// 		return 0, errors.NewBadRequestError("error couldn't find team with ID: " + requestData.TeamID)
// 	}
// 	if val.Password != requestData.TeamPassword {
// 		return 0, errors.NewWrongCredentialsError("incorrect team password")
// 	}
// 	delete(t.Teams, requestData.TeamID)
// 	return 1, nil
// }

// func (t *dbtest) AddMember(ctx context.Context, requestData accounts.request.AddMember) (bool, errors.APIErr) {
// 	var usr accounts.User
// 	for _, v := range t.Users {
// 		if v.Name == requestData.MemberName {
// 			usr = v
// 		}
// 	}
// 	if usr.ID == "" {
// 		return false, errors.NewBadRequestError("couldnt update user with name: " + requestData.MemberName)
// 	}
// 	mem, ok := t.Users[requestData.ID]
// 	team, found := t.Teams[requestData.TeamID]
// 	if !ok || !found || team.Members[requestData.ID] != accounts.RoleAdmin {
// 		return false, nil
// 	}
// 	team.Members[mem.ID] = requestData.Role
// 	return true, nil
// }

// func (t *dbtest) DeleteSelf(ctx context.Context, requestData accounts.request.DeleteSelf) (bool, errors.APIErr) {
// 	return true, nil
// }

// func (t *dbtest) DeleteMember(ctx context.Context, requestData accounts.request.DeleteMember) (bool, errors.APIErr) {
// 	return true, nil
// }

// func (t *dbtest) AddTeamCmd(ctx context.Context, requestData accounts.request.AddTeamCmd) (int, errors.APIErr) {
// 	return 1, nil
// }

// func (t *dbtest) DeleteTeamCmd(ctx context.Context, requestData accounts.request.DeleteTeamCmd, APIKey string) (int, errors.APIErr) {
// 	return 1, nil
// }
