package session

import "time"

//go:generate mockery --name=Service

const DefaultExpiry = 12 * time.Hour

type Service interface {
	Finder
	Saver
	Deleter
}
