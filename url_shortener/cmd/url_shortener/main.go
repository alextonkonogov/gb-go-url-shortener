package main

import (
	"context"
	_ "embed"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/api/handler"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/api/routeropenapi"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/db/pgstore"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/server"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/userrepo"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	ust, err := pgstore.NewUsers(os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	us := userrepo.NewURLs(ust)
	hs := handler.NewHandlers(us)
	h := routeropenapi.NewRouterOpenAPI(hs)
	srv := server.NewServer(":8000", h)

	srv.Start(us)
	log.Print("Start")

	<-ctx.Done()

	srv.Stop()
	cancel()
	ust.Close()

	log.Print("Exit")
}
