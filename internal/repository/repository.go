package repository

import (
	"context"
	"github.com/cr00z/gocqrs/internal/domain"
)

type RepositoryBasic interface {
	Close()
	Insert(ctx context.Context, meow domain.Meow) error
}

type RepositoryLister interface {
	RepositoryBasic
	List(ctx context.Context, skip, take uint64) ([]domain.Meow, error)
}

type RepositorySearcher interface {
	RepositoryBasic
	Search(ctx context.Context, query string, skip, take uint64) ([]domain.Meow, error)
}
