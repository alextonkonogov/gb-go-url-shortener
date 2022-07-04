package repoStatistics

import (
	"context"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/statistics"
)

// нужен только тут
type StatisticsStore interface {
	Read(ctx context.Context, s statistics.Statistics) (*statistics.Statistics, error)
}

type Statistics struct {
	StatisticsStore StatisticsStore
}

func NewStatistics(sstore StatisticsStore) *Statistics {
	return &Statistics{
		StatisticsStore: sstore,
	}
}

func (st *Statistics) Read(ctx context.Context, s statistics.Statistics) (*statistics.Statistics, error) {
	nst, err := st.StatisticsStore.Read(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("read statistics error: %w", err)
	}
	s.Long = nst.Long
	s.Short = nst.Short
	return &s, nil
}
