package api

type DatabaseAPI interface {
	List() error
	ListPages(database string) error
}

type database struct{}

var Database DatabaseAPI = database{}
