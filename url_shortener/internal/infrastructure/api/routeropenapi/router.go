package routeropenapi

import (
	"encoding/json"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/api/handler"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
)

type RouterOpenAPI struct {
	*chi.Mux
	hs  *handler.Handlers
	log *logrus.Logger
}

func NewRouterOpenAPI(hs *handler.Handlers, log *logrus.Logger) (*RouterOpenAPI, error) {
	r := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public", fileServer))

	ret := &RouterOpenAPI{
		hs:  hs,
		log: log,
	}

	r.Mount("/", Handler(ret))

	swg, err := GetSwagger()
	if err != nil {
		err = fmt.Errorf("swagger fail: %w", err)
		log.Error(err)
		return nil, err
	}

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		err = enc.Encode(swg)
		if err != nil {
			err = fmt.Errorf("swagger fail: %w", err)
			log.Error(err)
		}
	})

	ret.Mux = r
	return ret, nil
}
