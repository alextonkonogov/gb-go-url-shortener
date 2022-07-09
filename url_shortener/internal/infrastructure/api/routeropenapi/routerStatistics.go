package routeropenapi

import (
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/api/handler"
	"github.com/go-chi/render"
	"html/template"
	"net/http"
	"path/filepath"
)

type Statistics handler.Statistics

func (Statistics) Bind(r *http.Request) error {
	return nil
}

func (Statistics) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (rt *RouterOpenAPI) GetAAdmin(w http.ResponseWriter, r *http.Request, admin string) {
	page := filepath.Join("public", "html", "admin.html")
	common := filepath.Join("public", "html", "common.html")
	tmpl, err := template.ParseFiles(page, common)
	if err != nil {
		rt.log.WithError(fmt.Errorf("from route: %w", err)).Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "admin", nil)
	if err != nil {
		rt.log.WithError(fmt.Errorf("from route: %w", err)).Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *RouterOpenAPI) PostA(w http.ResponseWriter, r *http.Request) {
	s := Statistics{}
	if err := render.Bind(r, &s); err != nil {
		rt.log.WithError(fmt.Errorf("from route: %w", err)).Error(err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	st, err := rt.hs.ReadStatistics(r.Context(), handler.Statistics(s))
	if err != nil {
		rt.log.WithError(fmt.Errorf("from route: %w", err)).Error(err)
		render.Render(w, r, ErrInternalError(err))
		return
	}

	render.Render(w, r, Statistics(st))
}
