package service

import (
	"context"
	"errors"
	"time"

	"github.com/phihdn/golink/internal/repository"
)

const (
	DefaultCodeLength = 7
	// MaxGenerateRetries is how many codes CreateShortenLink tries before giving up.
	MaxGenerateRetries = 5
)

var (
	// ErrCodeExhausted is returned when every generated code was already taken.
	ErrCodeExhausted = errors.New("could not generate a unique code")
	// ErrLinkNotFound is returned by GetUrlFromCode when no link exists for the code.
	ErrLinkNotFound = errors.New("link not found")
)

type ShortenUrl interface {
	// CreateShortenLink stores url under a new unique code that expires after exp (0 = never) and returns the code.
	CreateShortenLink(ctx context.Context, url string, exp time.Duration) (string, error)
	// GetUrlFromCode returns the original url for code, or ErrLinkNotFound when the code is unknown or expired.
	GetUrlFromCode(ctx context.Context, code string) (string, error)
}

type shortenUrl struct {
	r       repository.URLStorage
	codeGen GenPass
}

func NewShortenUrl(r repository.URLStorage, codeGen GenPass) ShortenUrl {
	return &shortenUrl{r: r, codeGen: codeGen}
}

func (s *shortenUrl) CreateShortenLink(
	ctx context.Context,
	url string,
	exp time.Duration,
) (string, error) {
	for range MaxGenerateRetries {
		// generate code
		code, err := s.codeGen.GeneratePassword(DefaultCodeLength)
		if err != nil {
			return "", err
		}

		// store url theo code, retry with a new code when this one is taken
		err = s.r.StoreURL(ctx, code, url, exp)
		if errors.Is(err, repository.ErrKeyExisted) {
			continue
		}
		if err != nil {
			return "", err
		}

		// return code
		return code, nil
	}

	return "", ErrCodeExhausted
}

func (s *shortenUrl) GetUrlFromCode(ctx context.Context, code string) (string, error) {
	url, err := s.r.GetURL(ctx, code)
	if errors.Is(err, repository.ErrKeyNotFound) {
		return "", ErrLinkNotFound
	}
	if err != nil {
		return "", err
	}

	return url, nil
}
