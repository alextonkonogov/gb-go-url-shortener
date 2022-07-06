package repoStatistics

import (
	"context"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/statistics"
	"time"
)

// нужен только тут
type StatisticsStore interface {
	Create(ctx context.Context, URLID int64) error
	Read(ctx context.Context, s statistics.Statistics) (*statistics.Statistics, error)
	Update(ctx context.Context, s statistics.Statistics, URLID int64) (*statistics.Statistics, error)
}

type Statistics struct {
	StatisticsStore StatisticsStore
}

func NewStatistics(sstore StatisticsStore) *Statistics {
	return &Statistics{
		StatisticsStore: sstore,
	}
}

func (st *Statistics) Create(ctx context.Context, URLID int64) error {
	err := st.StatisticsStore.Create(ctx, URLID)
	if err != nil {
		return fmt.Errorf("create statistics error: %w", err)
	}
	return nil
}

func (st *Statistics) Read(ctx context.Context, s statistics.Statistics) (*statistics.Statistics, error) {
	nst, err := st.StatisticsStore.Read(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("read statistics error: %w", err)
	}
	s.IP = nst.IP
	s.Viewed = nst.Viewed
	s.Count = nst.Count
	s.Long = nst.Long
	s.Short = nst.Short

	return &s, nil
}

func (st *Statistics) Update(ctx context.Context, s statistics.Statistics, URLID int64) (*statistics.Statistics, error) {
	s.Viewed = time.Now().Format("2006-01-02T15:04:05.000Z")
	ns, err := st.StatisticsStore.Update(ctx, s, URLID)
	if err != nil {
		return nil, fmt.Errorf("update URL error: %w", err)
	}
	s.Count = ns.Count
	return &s, nil
}
