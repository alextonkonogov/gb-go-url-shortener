package handler

import (
	"errors"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repoStatistics"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repoURL"
)

type Handlers struct {
	ur *repoURL.URL
	st *repoStatistics.Statistics
}

func NewHandlers(ur *repoURL.URL, st *repoStatistics.Statistics) *Handlers {
	hs := &Handlers{
		ur: ur,
		st: st,
	}
	return hs
}

var ErrURLNotFound = errors.New("URL not found")
