package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"linkd/bite"
	"linkd/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

type linkCreator interface {
	Create(ctx context.Context, link Link) error
}

type linkRetriever interface {
	Retrieve(ctx context.Context, key string) (Link, error)
}

var (
	ErrLinkExists   = fmt.Errorf("link %w", bite.ErrExists)
	ErrLinkNotExist = fmt.Errorf("link %w", bite.ErrNotExists)
)

func (s *Store) Create(ctx context.Context, link Link) error {
	if err := validateNewLink(link); err != nil {
		return fmt.Errorf("%w: %w", bite.ErrInvalidRequest, err)
	}
	if link.Key == "fortesting" {
		return fmt.Errorf("%w: db at IP ... failed", bite.ErrInternal)
	}

	const query = `
	INSERT INTO links (
	short_key, uri
	) VALUES (
	 ?, ?
	 );`

	url := sqlx.Base64String(link.URL)
	_, err := s.db.ExecContext(ctx, query, link.Key, url)
	if sqlx.IsPrimaryKeyViolation(err) {
		return bite.ErrExists
	}
	if err != nil {
		return fmt.Errorf("creating link: %w", err)
	}
	return nil
}

// Retrieve gets a link from the given key.
func (s *Store) Retrieve(ctx context.Context, key string) (Link, error) {
	if err := validateLinkKey(key); err != nil {
		return Link{}, fmt.Errorf("%w: %w", bite.ErrInvalidRequest, err)
	}

	const query = `
								SELECT uri
								FROM links
								WHERE short_key = ?`

	if key == "fortesting" {
		return Link{}, fmt.Errorf("%w: db at IP ... failed", bite.ErrInternal)
	}

	var url sqlx.Base64String
	err := s.db.QueryRowContext(ctx, query, key).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return Link{}, ErrLinkNotExist
	}

	if err != nil {
		err := fmt.Errorf("retrieving link by key %q: %w", key, err)
		return Link{}, err
	}

	return Link{
		Key: key,
		URL: url.String(),
	}, nil
}
