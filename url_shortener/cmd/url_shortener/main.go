package main

import (
	"context"
	_ "embed"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/api/handler"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/api/routeropenapi"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/db/pgstore"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/server"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repoStatistics"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repoURL"
	"log"
	"os"
	"os/signal"
)

func main() {
	var pgStr = os.Getenv("PG_DSN")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	ust, err := pgstore.NewURL(pgStr)
	if err != nil {
		log.Fatal(err)
	}
	sst, err := pgstore.NewStatistics(pgStr)
	if err != nil {
		log.Fatal(err)
	}

	ur := repoURL.NewURL(ust)
	st := repoStatistics.NewStatistics(sst)
	hs := handler.NewHandlers(ur, st)
	h := routeropenapi.NewRouterOpenAPI(hs)
	srv := server.NewServer(":8000", h)

	srv.Start(ur, st)
	log.Print("Started")

	<-ctx.Done()

	srv.Stop()
	cancel()
	ust.Close()

	log.Print("Exit")
}
