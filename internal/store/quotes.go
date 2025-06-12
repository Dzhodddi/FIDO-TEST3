package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type Quote struct {
	ID     int64  `json:"id"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

const QueryTimeOut = time.Second * 10

type QuotesStore struct {
	db *sql.DB
}

func (s *QuotesStore) CreateQuote(ctx context.Context, quote *Quote) error {
	query := `INSERT INTO quotes (author, text) VALUES ($1, $2) RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, quote.Author, quote.Text).Scan(&quote.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *QuotesStore) GetByID(ctx context.Context, id int64) (*Quote, error) {
	query := `SELECT id, author, text FROM quotes WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	quote := &Quote{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&quote.ID, &quote.Author, &quote.Text)
	if err != nil {
		return nil, ErrNotFound
	}
	return quote, nil
}

func (s *QuotesStore) GetByAuthor(ctx context.Context, author string, paginatedQuery PaginatedQuery) ([]*Quote, error) {
	query := `SELECT id, author, text FROM quotes`
	args := []interface{}{}
	if author != "" {
		query += " WHERE author = $1 ORDER BY id LIMIT $2 OFFSET $3"
		args = append(args, author, paginatedQuery.Limit, paginatedQuery.Offset)
	} else {
		query += " ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, paginatedQuery.Limit, paginatedQuery.Offset)
	}

	log.Print(args)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var quotes []*Quote

	for rows.Next() {
		var quote Quote
		err = rows.Scan(&quote.ID, &quote.Author, &quote.Text)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, &quote)
	}
	return quotes, nil
}

func (s *QuotesStore) UpdateQuote(ctx context.Context, quote *Quote) error {
	query := `UPDATE quotes SET author = $1 , text = $2 WHERE ID = $3 RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, quote.Author, quote.Text, quote.ID).Scan(&quote.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return nil
		}
	}
	return nil
}

func (s *QuotesStore) DeleteQuote(ctx context.Context, id int64) error {
	query := `DELETE FROM quotes WHERE ID =  $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
