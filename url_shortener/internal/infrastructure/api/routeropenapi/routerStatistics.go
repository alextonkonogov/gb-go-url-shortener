package routeropenapi

import (
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tmpl.ExecuteTemplate(w, "admin", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *RouterOpenAPI) PostA(w http.ResponseWriter, r *http.Request) {
	s := Statistics{}
	if err := render.Bind(r, &s); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	u, err := rt.hs.ReadStatistics(r.Context(), handler.Statistics(s))
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Render(w, r, Statistics(u))
}
