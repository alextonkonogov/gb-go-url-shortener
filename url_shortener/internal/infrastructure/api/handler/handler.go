package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/url"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/userrepo"
)

type Handlers struct {
	us *userrepo.URLs
}

func NewHandlers(us *userrepo.URLs) *Handlers {
	r := &Handlers{
		us: us,
	}
	return r
}

type URL struct {
	ID    int    `json:"id"`
	Long  string `json:"long"`
	Short string `json:"short"`
}

func (rt *Handlers) CreateURL(ctx context.Context, u URL) (URL, error) {
	bu := url.URL{
		Long:  u.Long,
		Short: u.Short,
	}

	nbu, err := rt.us.Create(ctx, bu)
	if err != nil {
		return URL{}, fmt.Errorf("error when creating: %w", err)
	}

	return URL{
		ID:    nbu.ID,
		Long:  nbu.Long,
		Short: nbu.Short,
	}, nil
}

var ErrURLNotFound = errors.New("URL not found")
