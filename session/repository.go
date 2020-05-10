package session

//go:generate mockery --name=Repository

type Repository interface {
	Finder
	Saver
	Deleter
}
