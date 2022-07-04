package repoURL

import (
	"context"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/url"
	"github.com/dchest/uniuri"
	"time"
)

// нужен только тут
type URLStore interface {
	Create(ctx context.Context, u url.URL) (*int64, error)
	Read(ctx context.Context, u url.URL) (*string, error)
}

type URL struct {
	URLStore URLStore
}

func NewURL(ustore URLStore) *URL {
	return &URL{
		URLStore: ustore,
	}
}

func (us *URL) Create(ctx context.Context, u url.URL) (*url.URL, error) {
	u.Created = time.Now().Format("2006-01-02T15:04:05.000Z")
	u.Short = uniuri.New()
	u.Admin = uniuri.New()
	id, err := us.URLStore.Create(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("create URL error: %w", err)
	}
	u.ID = *id
	return &u, nil
}

func (us *URL) Read(ctx context.Context, u url.URL) (*url.URL, error) {
	long, err := us.URLStore.Read(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("create URL error: %w", err)
	}
	u.Long = *long
	return &u, nil
}
