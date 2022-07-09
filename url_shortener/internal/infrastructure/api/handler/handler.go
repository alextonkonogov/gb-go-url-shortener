package handler

import (
	"errors"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repoStatistics"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repoURL"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	ur  *repoURL.URL
	st  *repoStatistics.Statistics
	log *logrus.Logger
}

func NewHandlers(ur *repoURL.URL, st *repoStatistics.Statistics, log *logrus.Logger) *Handlers {
	hs := &Handlers{
		ur:  ur,
		st:  st,
		log: log,
	}
	return hs
}

var ErrURLNotFound = errors.New("URL not found")
