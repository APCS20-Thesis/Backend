package store

type Querier interface {
}

var _ Querier = (*Queries)(nil)
