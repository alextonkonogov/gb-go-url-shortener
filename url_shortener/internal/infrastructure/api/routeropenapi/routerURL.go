package routeropenapi

import (
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/statistics"
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

func (rt *RouterOpenAPI) GetErr(w http.ResponseWriter, r *http.Request) {
	page := filepath.Join("public", "html", "err.html")
	common := filepath.Join("public", "html", "common.html")
	tmpl, err := template.ParseFiles(page, common)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tmpl.ExecuteTemplate(w, "err", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *RouterOpenAPI) PostSCreate(w http.ResponseWriter, r *http.Request) {
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

	err = rt.hs.CreateStatistics(r.Context(), u.ID)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Render(w, r, URL(u))
}

func (rt *RouterOpenAPI) GetSShort(w http.ResponseWriter, r *http.Request, short string) {
	u := URL{}
	u.Short = short

	nu, err := rt.hs.ReadURL(r.Context(), handler.URL(u))
	if err != nil {
		http.Redirect(w, r, "/err", http.StatusSeeOther)
		return
	}

	st := statistics.Statistics{}
	st.IP = r.RemoteAddr
	_, err = rt.hs.UpdateStatistics(r.Context(), handler.Statistics(st), nu.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, nu.Long, http.StatusSeeOther)
}
