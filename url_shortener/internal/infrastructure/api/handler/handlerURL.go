package handler

import (
	"context"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/url"
)

type URL struct {
	ID      int64  `json:"id"`
	Created string `json:"created"`
	Long    string `json:"long"`
	Short   string `json:"short"`
	Admin   string `json:"admin"`
}

func (h *Handlers) CreateURL(ctx context.Context, u URL) (URL, error) {
	bu := url.URL{
		Long: u.Long,
	}

	nbu, err := h.ur.Create(ctx, bu)
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

func (h *Handlers) ReadURL(ctx context.Context, u URL) (URL, error) {
	bu := url.URL{
		Short: u.Short,
	}

	nbu, err := h.ur.Read(ctx, bu)
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
