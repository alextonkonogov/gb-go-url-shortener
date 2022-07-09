package main

import (
	"context"
	_ "embed"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/api/handler"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/api/routeropenapi"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/db/pgstore"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/log"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/server"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repostatistics"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repourl"

	"os"
	"os/signal"
)

var appIP = os.Getenv("APP_IP")
var appPort = os.Getenv("APP_PORT")
var pgStr = os.Getenv("DB_CONNECTION_STRING")

func main() {
	l := log.NewLogWithConfuguration()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	ust, err := pgstore.NewURL(pgStr, l)
	if err != nil {
		l.WithError(err).Fatal()
	}
	sst, err := pgstore.NewStatistics(pgStr, l)
	if err != nil {
		l.WithError(err).Fatal()
	}

	ur := repourl.NewURL(ust, l)
	st := repostatistics.NewStatistics(sst, l)
	hs := handler.NewHandlers(ur, st, l)
	h, err := routeropenapi.NewRouterOpenAPI(hs, l)
	if err != nil {
		l.WithError(err).Fatal()
	}
	srv := server.NewServer(appIP+":"+appPort, h, l)

	srv.Start(ur, st)
	l.Info("started")

	<-ctx.Done()

	srv.Stop()
	cancel()
	ust.Close()

	l.Info("exited")
}
