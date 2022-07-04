package server

import (
	"context"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repoStatistics"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repoURL"
	"net/http"
	"time"
)

type Server struct {
	srv http.Server
	ur  *repoURL.URL
	st  *repoStatistics.Statistics
}

func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

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
	s.srv.Shutdown(ctx)
	cancel()
}

func (s *Server) Start(ur *repoURL.URL, st *repoStatistics.Statistics) {
	s.ur = ur
	s.st = st
	// TODO: migrations
	go s.srv.ListenAndServe()
}
