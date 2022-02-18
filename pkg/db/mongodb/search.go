package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/conalli/bookshelf-backend/pkg/db"
)

// Search finds the user in the DB and returns the url for a given command.
func (m *Mongo) Search(ctx context.Context, APIKey, cmd string) (string, error) {
	reqCtx, cancelFunc := db.ReqContextWithTimeout(ctx)
	defer cancelFunc()
	m.Initialize()
	err := m.client.Connect(reqCtx)
	if err != nil {
		log.Println("couldjnt connect to db on search")
	}
	defer m.client.Disconnect(reqCtx)
	collection := m.db.Collection(CollectionUsers)
	res := GetByKey(ctx, collection, "APIKey", APIKey)
	currUser, err := DecodeUser(res)
	defaultSearch := fmt.Sprintf("http://www.google.com/search?q=%s", cmd)
	if err != nil {
		return defaultSearch, err
	}
	url, ok := currUser.Bookmarks[cmd]
	if !ok {
		return defaultSearch, err
	}
	return url, nil
}
