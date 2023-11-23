package repository

import (
	"context"

	"github.com/cr00z/gocqrs/internal/domain"
)

type RepositoryBasic interface {
	Close()
	Insert(ctx context.Context, msg domain.Message) error
}

type RepositoryLister interface {
	RepositoryBasic
	List(ctx context.Context, skip, take uint64) ([]domain.Message, error)
}

type RepositorySearcher interface {
	RepositoryBasic
	Search(ctx context.Context, query string, skip, take uint64) ([]domain.Message, error)
}
