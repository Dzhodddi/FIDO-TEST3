package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("Record not found")
)

type Storage struct {
	Quotes interface {
		CreateQuote(ctx context.Context, quote *Quote) error
		GetByID(ctx context.Context, id int64) (*Quote, error)
		GetByAuthor(ctx context.Context, author string, query PaginatedQuery) ([]*Quote, error)
		UpdateQuote(ctx context.Context, quote *Quote) error
		DeleteQuote(ctx context.Context, id int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Quotes: &QuotesStore{db},
	}
}
