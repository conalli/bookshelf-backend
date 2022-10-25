package db

import (
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
	"github.com/conalli/bookshelf-backend/pkg/services/search"
)

// Storage represents the storage for all services.
type Storage interface {
	accounts.UserRepository
	bookmarks.Repository
	search.Repository
}

// Cache represents the cache for all services.
type Cache interface {
	accounts.UserCache
	search.Cache
}
