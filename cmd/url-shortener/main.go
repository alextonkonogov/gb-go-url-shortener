package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/alextonkonogov/gb-go-url-shortener/internal/application"
	"github.com/alextonkonogov/gb-go-url-shortener/internal/storage"
)

func main() {
	ctx := context.Background()

	dbpool, err := storage.InitDBConn(ctx)
	if err != nil {
		log.Panic(fmt.Errorf("%w failed to init DB connection", err))
	}
	defer dbpool.Close()

	err = storage.InitTables(ctx, dbpool)
	if err != nil {
		log.Panic(fmt.Errorf("%w failed to init DB tables", err))
	}

	app := application.NewApp(ctx, dbpool)
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	app.Routes(r)

	if err = http.ListenAndServe("0.0.0.0:8282", r); err != nil {
		log.Panic(fmt.Errorf("%w failed to listen and serve", err))
	}
}
