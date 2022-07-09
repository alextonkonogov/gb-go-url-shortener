package repostatistics

import (
	"context"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/statistics"
	"github.com/sirupsen/logrus"
	"time"
)

// нужен только тут
type StatisticsStore interface {
	Create(ctx context.Context, urlId int64) error
	Read(ctx context.Context, s statistics.Statistics) (*statistics.Statistics, error)
	Update(ctx context.Context, s statistics.Statistics, urlId int64) (*statistics.Statistics, error)
}

type Statistics struct {
	StatisticsStore StatisticsStore
	log             *logrus.Logger
}

func NewStatistics(sstore StatisticsStore, log *logrus.Logger) *Statistics {
	return &Statistics{
		StatisticsStore: sstore,
		log:             log,
	}
}

func (st *Statistics) Create(ctx context.Context, urlId int64) error {
	err := st.StatisticsStore.Create(ctx, urlId)
	if err != nil {
		err = fmt.Errorf("create statistics error: %w", err)
		st.log.Error(err)
		return err
	}
	return nil
}

func (st *Statistics) Read(ctx context.Context, s statistics.Statistics) (*statistics.Statistics, error) {
	nst, err := st.StatisticsStore.Read(ctx, s)
	if err != nil {
		err = fmt.Errorf("read statistics error: %w", err)
		st.log.Error(err)
		return nil, err
	}
	s.IP = nst.IP
	s.Viewed = nst.Viewed
	s.Count = nst.Count
	s.Long = nst.Long
	s.Short = nst.Short

	return &s, nil
}

func (st *Statistics) Update(ctx context.Context, s statistics.Statistics, urlId int64) (*statistics.Statistics, error) {
	s.Viewed = time.Now().Format("2006-01-02T15:04:05.000Z")
	ns, err := st.StatisticsStore.Update(ctx, s, urlId)
	if err != nil {
		err = fmt.Errorf("update statistics error: %w", err)
		st.log.Error(err)
		return nil, err
	}
	s.Count = ns.Count
	return &s, nil
}
