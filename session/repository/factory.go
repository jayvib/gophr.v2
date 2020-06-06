package repository

import (
	"github.com/patrickmn/go-cache"
	"gophr.v2/config"
	"gophr.v2/driver/redis"
	"gophr.v2/session"
	"gophr.v2/session/repository/file"
	"gophr.v2/session/repository/gocache"
	redisrepo "gophr.v2/session/repository/redis"
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
	// RedisRepo describes the redis implementation of the repository.
	RedisRepo
)

// Get is a factory function that accepts rt repository type
// and return the implementation of the session Repository interface.
func Get(conf *config.Config, rt RepoType) session.Repository {
	switch rt {
	case FileRepo:
		return file.New(file.DefaultFilename)
	case GoCacheRepo:
		c := cache.New(gocache.DefaultExpirationTime, gocache.DefaultExpirationTime)
		return gocache.New(c)
	case RedisRepo:
		conn := redis.New(conf)
		return redisrepo.New(conn)
	default:
		panic("unknown repository type implementation")
	}
}
