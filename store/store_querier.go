package store

type StoreQuerier interface {
	Querier
	Ping() error
	Close() error
	//Some Other Query func here
}
