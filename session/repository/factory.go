package repository

import (
	"github.com/patrickmn/go-cache"
	"gophr.v2/session"
	"gophr.v2/session/repository/file"
	"gophr.v2/session/repository/gocache"
)

// factory design pattern. Its purpose is to abstract the user from the knowledge
// of the struct he needs to achieve for a specific purpose.

// RepoType describes the repository type
type RepoType int

const (
	// FileRepo describes the file repository implementation type
	FileRepo RepoType = iota
	// GoCacheRepo describes the gocache repository implementation type
	GoCacheRepo
)

// GetRepository is a factory function that accepts rt repository type
// and return the implementation of the session Repository interface.
func GetRepository(rt RepoType) session.Repository {
	switch rt {
	case FileRepo:
		return file.New(file.DefaultFilename)
	case GoCacheRepo:
		c := cache.New(gocache.DefaultExpirationTime, gocache.DefaultExpirationTime)
		return gocache.New(c)
	default:
		panic("unknown repository type implementation")
	}
}
