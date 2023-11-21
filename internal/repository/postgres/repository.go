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

func (r PostgresRepository) List(ctx context.Context, skip, take uint64) ([]domain.Message, error) {
	rows, err := r.db.Query("SELECT id, body, created_at FROM messages ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		msg := domain.Message{}
		err = rows.Scan(&msg.ID, msg.Body, msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r PostgresRepository) Insert(ctx context.Context, msg domain.Message) error {
	_, err := r.db.Exec("INSERT INTO messages(id, body, created_at) VALUES($1, $2, $3)",
		msg.ID, msg.Body, msg.CreatedAt)
	return err
}
