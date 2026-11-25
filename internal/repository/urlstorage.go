package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrKeyExisted is returned by StoreURL when the key is already taken.
	ErrKeyExisted = errors.New("key already exists")
	// ErrKeyNotFound is returned by GetURL when the key does not exist or has expired.
	ErrKeyNotFound = errors.New("key not found")
)

type URLStorage interface {
	// StoreURL saves key -> url only when the key is not taken yet, otherwise it returns ErrKeyExisted.
	// A zero exp means the url never expires.
	StoreURL(ctx context.Context, key, url string, exp time.Duration) error
	// GetURL returns the url saved under key, or ErrKeyNotFound when there is none.
	GetURL(ctx context.Context, key string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

func NewURLStorage(c *redis.Client) URLStorage {
	return &urlStorage{c: c}
}

func (s *urlStorage) StoreURL(ctx context.Context, key, url string, exp time.Duration) error {
	// SET NX checks and writes in one atomic step, so two requests can never claim the same key.
	stored, err := s.c.SetNX(ctx, key, url, exp).Result()
	if err != nil {
		return err
	}
	if !stored {
		return ErrKeyExisted
	}

	return nil
}

func (s *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	url, err := s.c.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrKeyNotFound
	}
	if err != nil {
		return "", err
	}

	return url, nil
}
