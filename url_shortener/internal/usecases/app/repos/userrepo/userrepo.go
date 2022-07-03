package userrepo

import (
	"context"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/url"
)

// нужен только тут
type URLStore interface {
	Create(ctx context.Context, u url.URL) (*int, error)
}

type URLs struct {
	ustore URLStore
}

func NewURLs(ustore URLStore) *URLs {
	return &URLs{
		ustore: ustore,
	}
}

func (us *URLs) Create(ctx context.Context, u url.URL) (*url.URL, error) {
	return &u, nil
}
