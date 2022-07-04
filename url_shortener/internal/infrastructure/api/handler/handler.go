package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/url"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/URLrepo"
)

type Handlers struct {
	us *URLrepo.URLs
}

func NewHandlers(us *URLrepo.URLs) *Handlers {
	r := &Handlers{
		us: us,
	}
	return r
}

type URL struct {
	ID      int64  `json:"id"`
	Created string `json:"created"`
	Long    string `json:"long"`
	Short   string `json:"short"`
	Admin   string `json:"admin"`
}

func (rt *Handlers) CreateURL(ctx context.Context, u URL) (URL, error) {
	bu := url.URL{
		Long: u.Long,
	}

	nbu, err := rt.us.Create(ctx, bu)
	if err != nil {
		return URL{}, fmt.Errorf("error when creating: %w", err)
	}

	return URL{
		ID:      nbu.ID,
		Created: nbu.Created,
		Long:    nbu.Long,
		Short:   nbu.Short,
		Admin:   nbu.Admin,
	}, nil
}

func (rt *Handlers) ReadURL(ctx context.Context, u URL) (URL, error) {
	bu := url.URL{
		Short: u.Short,
	}

	nbu, err := rt.us.Read(ctx, bu)
	if err != nil {
		return URL{}, fmt.Errorf("error when reading: %w", err)
	}

	return URL{
		ID:      nbu.ID,
		Created: nbu.Created,
		Long:    nbu.Long,
		Short:   nbu.Short,
		Admin:   nbu.Admin,
	}, nil
}

var ErrURLNotFound = errors.New("URL not found")
