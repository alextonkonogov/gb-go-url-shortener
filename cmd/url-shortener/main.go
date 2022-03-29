package main

import (
	"context"
	"github.com/alextonkonogov/gb-go-url-shortener/internal/storage"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/alextonkonogov/gb-go-url-shortener/internal/application"
)

func main() {
	ctx := context.Background()

	dbpool, err := storage.InitDBConn(ctx)
	if err != nil {
		log.Fatalf("%w failed to init DB connection", err)
	}
	defer dbpool.Close()

	err = storage.InitTables(ctx, dbpool)
	if err != nil {
		log.Fatalf("%w failed to init DB tables", err)
	}

	app := application.NewApp(ctx, dbpool)
	r := httprouter.New()
	app.Routes(r)

	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: r}
	if err = srv.ListenAndServe(); err != nil {
		log.Fatalf("%w failed to listen and serve", err)
	}
}
