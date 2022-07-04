package handler

import (
	"context"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/statistics"
)

type Statistics struct {
	IP      string `json:"ip"`
	Visited string `json:"visited"`
	Count   int    `json:"count"`
	Long    string `json:"long"`
	Short   string `json:"short"`
	Admin   string `json:"admin"`
}

func (h *Handlers) ReadStatistics(ctx context.Context, st Statistics) (Statistics, error) {
	bu := statistics.Statistics{
		Admin: st.Admin,
	}

	nbu, err := h.st.Read(ctx, bu)
	if err != nil {
		return Statistics{}, fmt.Errorf("error when reading: %w", err)
	}

	return Statistics{
		Count:   nbu.Count,
		Visited: nbu.Visited,
		IP:      nbu.IP,
		Long:    nbu.Long,
		Short:   nbu.Short,
	}, nil
}
