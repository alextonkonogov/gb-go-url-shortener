package routeropenapi

import (
	"encoding/json"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/api/handler"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type RouterOpenAPI struct {
	*chi.Mux
	hs *handler.Handlers
}

func NewRouterOpenAPI(hs *handler.Handlers) *RouterOpenAPI {
	r := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public", fileServer))

	ret := &RouterOpenAPI{
		hs: hs,
	}

	r.Mount("/", Handler(ret))

	swg, err := GetSwagger()
	if err != nil {
		log.Fatal("swagger fail")
	}

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		_ = enc.Encode(swg)
	})

	ret.Mux = r
	return ret
}
