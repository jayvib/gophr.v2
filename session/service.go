package session

//go:generate mockery --name=Service

type Service interface {
	Finder
	Saver
	Deleter
}
