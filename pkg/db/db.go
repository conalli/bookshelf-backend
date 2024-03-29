package db

import (
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
	"github.com/conalli/bookshelf-backend/pkg/services/search"
)

// Storage represents the storage for all services.
type Storage interface {
	auth.Repository
	accounts.UserRepository
	bookmarks.Repository
	search.Repository
}

// Cache represents the cache for all services.
type Cache interface {
	auth.Cache
	accounts.UserCache
	search.Cache
}
