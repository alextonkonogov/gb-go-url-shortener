package handler

import (
	"context"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/statistics"
)

type Statistics struct {
	IP     string `json:"ip"`
	Viewed string `json:"viewed"`
	Count  int64  `json:"count"`
	Long   string `json:"long"`
	Short  string `json:"short"`
	Admin  string `json:"admin"`
}

func (h *Handlers) CreateStatistics(ctx context.Context, urlId int64) error {
	err := h.st.Create(ctx, urlId)
	if err != nil {
		err = fmt.Errorf("error when creating: %w", err)
		h.log.Error(err)
		return err
	}

	return nil
}

func (h *Handlers) ReadStatistics(ctx context.Context, st Statistics) (Statistics, error) {
	bu := statistics.Statistics{
		Admin: st.Admin,
	}

	nbu, err := h.st.Read(ctx, bu)
	if err != nil {
		err = fmt.Errorf("error when reading: %w", err)
		h.log.Error(err)
		return Statistics{}, err
	}

	return Statistics{
		Count:  nbu.Count,
		Viewed: nbu.Viewed,
		IP:     nbu.IP,
		Long:   nbu.Long,
		Short:  nbu.Short,
	}, nil
}

func (h *Handlers) UpdateStatistics(ctx context.Context, st Statistics, urlId int64) (Statistics, error) {
	bu := statistics.Statistics{
		Admin: st.Admin,
		IP:    st.IP,
	}

	nbu, err := h.st.Update(ctx, bu, urlId)
	if err != nil {
		err = fmt.Errorf("error when updating: %w", err)
		h.log.Error(err)
		return Statistics{}, err
	}

	return Statistics{
		Count:  nbu.Count,
		Viewed: nbu.Viewed,
		IP:     nbu.IP,
		Long:   nbu.Long,
		Short:  nbu.Short,
	}, nil
}
