package repoURL

import (
	"context"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/url"
	"github.com/dchest/uniuri"
	"github.com/sirupsen/logrus"
	"time"
)

// нужен только тут
type URLStore interface {
	Create(ctx context.Context, u url.URL) (*int64, error)
	Read(ctx context.Context, u url.URL) (*url.URL, error)
}

type URL struct {
	URLStore URLStore
	log      *logrus.Logger
}

func NewURL(ustore URLStore, log *logrus.Logger) *URL {
	return &URL{
		URLStore: ustore,
		log:      log,
	}
}

func (ur *URL) Create(ctx context.Context, u url.URL) (*url.URL, error) {
	u.Created = time.Now().Format("2006-01-02T15:04:05.000Z")
	u.Short = uniuri.New()
	u.Admin = uniuri.New()
	id, err := ur.URLStore.Create(ctx, u)
	if err != nil {
		err = fmt.Errorf("create URL error: %w", err)
		ur.log.Error(err)
		return nil, err
	}
	u.ID = *id
	return &u, nil
}

func (ur *URL) Read(ctx context.Context, u url.URL) (*url.URL, error) {
	nu, err := ur.URLStore.Read(ctx, u)
	if err != nil {
		err = fmt.Errorf("read URL error: %w", err)
		ur.log.Error(err)
		return nil, err
	}
	u.ID = nu.ID
	u.Long = nu.Long
	return &u, nil
}
