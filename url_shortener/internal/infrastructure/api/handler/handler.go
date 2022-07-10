package handler

import (
	"errors"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repostatistics"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repourl"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	ur  *repourl.URL
	st  *repostatistics.Statistics
	log *logrus.Logger
}

func NewHandlers(ur *repourl.URL, st *repostatistics.Statistics, log *logrus.Logger) *Handlers {
	hs := &Handlers{
		ur:  ur,
		st:  st,
		log: log,
	}
	return hs
}

var ErrURLNotFound = errors.New("URL not found")
