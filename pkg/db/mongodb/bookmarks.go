package mongodb

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllBookmarks gets all a users bookmarks from the db.
func (m *Mongo) GetAllBookmarks(ctx context.Context, APIKey string) ([]bookmarks.Bookmark, apierr.Error) {
	collection := m.db.Collection(CollectionBookmarks)
	filter := bson.D{primitive.E{Key: "api_key", Value: APIKey}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		m.log.Errorf("could not find all bookmarks by APIKey: %v", err)
		return nil, apierr.NewInternalServerError()
	}
	var bookmarks []bookmarks.Bookmark
	err = cursor.All(ctx, &bookmarks)
	if err != nil {
		m.log.Errorf("could not get bookmarks from db cursor: %v", err)
		return nil, apierr.NewInternalServerError()
	}
	return bookmarks, nil
}

// GetBookmarksFolder gets all a users bookmarks from the db.
func (m *Mongo) GetBookmarksFolder(ctx context.Context, path, APIKey string) ([]bookmarks.Bookmark, apierr.Error) {
	collection := m.db.Collection(CollectionBookmarks)
	filter := bson.D{
		bson.E{Key: "api_key", Value: APIKey},
		bson.E{
			Key: "$or", Value: bson.A{
				bson.M{
					"path": bson.D{
						bson.E{Key: "$regex", Value: primitive.Regex{Pattern: path, Options: "i"}},
					},
				},
				bson.M{
					"is_folder": true,
					"name": bson.D{
						bson.E{Key: "$regex", Value: primitive.Regex{Pattern: path, Options: "i"}},
					},
				},
			},
		},
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		m.log.Errorf("could not find all bookmarks by APIKey: %v", err)
		return nil, apierr.NewInternalServerError()
	}
	var bookmarks []bookmarks.Bookmark
	err = cursor.All(ctx, &bookmarks)
	if err != nil {
		m.log.Errorf("could not get bookmarks from db cursor: %v", err)
		return nil, apierr.NewInternalServerError()
	}
	return bookmarks, nil
}

// AddBookmark adds a new bookmark for a given user.
func (m *Mongo) AddBookmark(ctx context.Context, requestData request.AddBookmark, APIKey string) (int, apierr.Error) {
	collection := m.db.Collection(CollectionBookmarks)
	data := bookmarks.Bookmark{
		APIKey:   APIKey,
		Name:     requestData.Name,
		Path:     requestData.Path,
		URL:      requestData.URL,
		IsFolder: requestData.IsFolder,
	}
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		m.log.Errorf("couldn't insert bookmark: %v", err)
		return 0, apierr.NewInternalServerError()
	}
	return 1, nil
}

func (m *Mongo) AddManyBookmarks(ctx context.Context, bookmarks []bookmarks.Bookmark) (int, apierr.Error) {
	collection := m.db.Collection(CollectionBookmarks)
	data := make([]interface{}, len(bookmarks))
	for i := range bookmarks {
		data[i] = bookmarks[i]
	}
	res, err := collection.InsertMany(ctx, data)
	if err != nil {
		m.log.Errorf("could not insert many bookmarks into db - %v", err)
		return 0, apierr.NewInternalServerError()
	}
	m.log.Infof("inserted %d bookmarks into db", len(res.InsertedIDs))
	return len(res.InsertedIDs), nil
}

// DeleteBookmark removes a bookmark for a given user.
func (m *Mongo) DeleteBookmark(ctx context.Context, bookmarkID, APIKey string) (int, apierr.Error) {
	collection := m.db.Collection(CollectionBookmarks)
	oid, err := primitive.ObjectIDFromHex(bookmarkID)
	if err != nil {
		m.log.Error("could not get ObjectID from Hex")
		return 0, apierr.NewBadRequestError("invalid bookmark id")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: oid}}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		m.log.Errorf("couldn't remove cmd from user: %v", err)
		return 0, apierr.NewInternalServerError()
	}
	return int(result.DeletedCount), nil
}
