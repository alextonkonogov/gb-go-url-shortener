package defmux

import (
	"encoding/json"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/api/handler"
	"net/http"
)

type Router struct {
	*http.ServeMux
	hs *handler.Handlers
}

func NewRouter(hs *handler.Handlers) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		hs:       hs,
	}

	r.Handle("/create", http.HandlerFunc(r.CreateUser))
	return r
}

func (rt *Router) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	u := handler.User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	nbu, err := rt.hs.CreateURL(r.Context(), u)
	if err != nil {
		http.Error(w, "error when creating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(nbu)
}
