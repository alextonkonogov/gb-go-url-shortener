package server

import (
	"context"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repostatistics"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repourl"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	srv http.Server
	ur  *repourl.URL
	st  *repostatistics.Statistics
	log *logrus.Logger
}

func NewServer(addr string, h http.Handler, log *logrus.Logger) *Server {
	s := &Server{
		log: log,
	}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	return s
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.WithError(err).Fatal()
	}
	cancel()
}

func (s *Server) Start(ur *repourl.URL, st *repostatistics.Statistics) {
	s.ur = ur
	s.st = st
	// TODO: migrations
	go func() {
		_ = s.srv.ListenAndServe()
	}()
}
