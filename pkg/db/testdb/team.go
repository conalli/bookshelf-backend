package testdb

// import (
// 	"context"

// 	"github.com/conalli/bookshelf-backend/pkg/accounts"
// 	"github.com/conalli/bookshelf-backend/pkg/errors"
// )

// func (t *testdb) NewTeam(ctx context.Context, requestData accounts.NewTeamRequest) (string, errors.ApiErr) {
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

// func (t *testdb) DeleteTeam(ctx context.Context, requestData accounts.DelTeamRequest) (int, errors.ApiErr) {
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

// func (t *testdb) AddMember(ctx context.Context, requestData accounts.AddMemberRequest) (bool, errors.ApiErr) {
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

// func (t *testdb) DeleteSelf(ctx context.Context, requestData accounts.DelSelfRequest) (bool, errors.ApiErr) {
// 	return true, nil
// }

// func (t *testdb) DeleteMember(ctx context.Context, requestData accounts.DelMemberRequest) (bool, errors.ApiErr) {
// 	return true, nil
// }

// func (t *testdb) AddTeamCmd(ctx context.Context, requestData accounts.AddTeamCmdRequest) (int, errors.ApiErr) {
// 	return 1, nil
// }

// func (t *testdb) DeleteTeamCmd(ctx context.Context, requestData accounts.DelTeamCmdRequest, APIKey string) (int, errors.ApiErr) {
// 	return 1, nil
// }
