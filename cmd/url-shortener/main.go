package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"

	"github.com/alextonkonogov/gb-go-url-shortener/internal/application"
	"github.com/alextonkonogov/gb-go-url-shortener/internal/storage"
)

var APP_IP = os.Getenv("APP_IP")
var APP_PORT = os.Getenv("APP_PORT")
var DB_CONNECTION_STRING = os.Getenv("DB_CONNECTION_STRING")

func main() {

	ctx := context.Background()

	dbpool, err := storage.InitDBConn(ctx, DB_CONNECTION_STRING)
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

	if err = http.ListenAndServe(APP_IP+":"+APP_PORT, r); err != nil {
		log.Panic(fmt.Errorf("%w failed to listen and serve", err))
	}
}
