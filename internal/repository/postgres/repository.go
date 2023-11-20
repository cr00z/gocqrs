package postgres

import (
	"context"
	"database/sql"

	"github.com/cr00z/gocqrs/internal/domain"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

func (r PostgresRepository) Close() {
	r.db.Close()
}

func (r PostgresRepository) List(ctx context.Context, skip, take uint64) ([]domain.Meow, error) {
	rows, err := r.db.Query("SELECT id, body, created_at FROM meows ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meows []domain.Meow
	for rows.Next() {
		meow := domain.Meow{}
		err = rows.Scan(&meow.ID, meow.Body, meow.CreatedAt)
		if err != nil {
			return nil, err
		}
		meows = append(meows, meow)
	}

	return meows, nil
}

func (r PostgresRepository) Insert(ctx context.Context, meow domain.Meow) error {
	_, err := r.db.Exec("INSERT INTO meows(id, body, created_at) VALUES($1, $2, $3)",
		meow.ID, meow.Body, meow.CreatedAt)
	return err
}
