package routeropenapi

import (
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/infrastructure/api/handler"
	"github.com/go-chi/render"
	"html/template"
	"net/http"
	"path/filepath"
)

type URL handler.URL

func (URL) Bind(r *http.Request) error {
	return nil
}

func (URL) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (rt *RouterOpenAPI) Get(w http.ResponseWriter, r *http.Request) {
	index := filepath.Join("public", "html", "index.html")
	common := filepath.Join("public", "html", "common.html")
	tmpl, err := template.ParseFiles(index, common)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *RouterOpenAPI) PostCreate(w http.ResponseWriter, r *http.Request) {
	ru := URL{}
	if err := render.Bind(r, &ru); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	u, err := rt.hs.CreateURL(r.Context(), handler.URL(ru))
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Render(w, r, URL(u))
}

func (rt *RouterOpenAPI) GetSShort(w http.ResponseWriter, r *http.Request, short string) {
	ru := URL{}
	ru.Short = short

	u, err := rt.hs.ReadURL(r.Context(), handler.URL(ru))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, u.Long, http.StatusSeeOther)
}
