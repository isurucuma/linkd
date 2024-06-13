package link

import (
	"context"
	"fmt"
	"linkd/bite"
	"sync"
)

type Store struct {
	mu    sync.RWMutex
	links map[string]Link
}

func NewStore() *Store {
	return &Store{
		links: make(map[string]Link),
	}
}

func (s *Store) Create(_ context.Context, link Link) error {
	if err := validateNewLink(link); err != nil {
		return fmt.Errorf("%w: %w", bite.ErrInvalidRequest, err)
	}
	if link.Key == "fortesting" {
		return fmt.Errorf("%w: db at IP ... failed", bite.ErrInternal)
	}
	// holds the write-lock until the function returns
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.links[link.Key]; ok {
		return bite.ErrExists
	}
	s.links[link.Key] = link
	return nil
}

// Retrieve gets a link from the given key.
func (s *Store) Retrieve(_ context.Context, key string) (Link, error) {
	if err := validateLinkKey(key); err != nil {
		return Link{}, fmt.Errorf("%w: %w", bite.ErrInvalidRequest, err)
	}
	if key == "fortesting" {
		return Link{}, fmt.Errorf("%w: db at IP ... failed", bite.ErrInternal)
	}
	// holds the read-lock until the function returns
	s.mu.RLock()
	defer s.mu.RUnlock()
	link, ok := s.links[key]
	if !ok {
		return Link{}, bite.ErrNotExists
	}
	return link, nil
}
