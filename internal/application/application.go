package application

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/alextonkonogov/gb-go-url-shortener/internal/repository"
	"github.com/dchest/uniuri"
	"github.com/jackc/pgx/v4/pgxpool"
)

const URLregexp = `(?m)https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&\/\/=]*)`

type app struct {
	ctx    context.Context
	dbpool *pgxpool.Pool
	repo   *repository.Repository
}

func (a app) Routes(r *chi.Mux) {
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/public/*", http.StripPrefix("/public", fileServer))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		a.IndexPage(w, nil)
	})

	r.Post("/short", a.ShortURL)
	r.Get("/s/{ID}/{code}", a.LongToShort)
	r.Get("/a/{ID}/{code}", a.AdminsPage)
}

func (a app) IndexPage(rw http.ResponseWriter, data interface{}) {
	index := filepath.Join("public", "html", "index.html")
	common := filepath.Join("public", "html", "common.html")
	tmpl, err := template.ParseFiles(index, common)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = tmpl.ExecuteTemplate(rw, "index", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) ShortURL(w http.ResponseWriter, r *http.Request) {
	var longURL, shortURL, adminURL string

	longURL = r.FormValue("longURL")

	data := struct {
		Err     bool
		Content []struct {
			Title string
			Link  string
		}
	}{}

	if !regexp.MustCompile(URLregexp).MatchString(longURL) {
		data.Err = true
		a.IndexPage(w, data)
		return
	}

	shortURL, adminURL = uniuri.New(), uniuri.New()

	shortURLID, err := a.repo.NewShortURL(a.ctx, a.dbpool, shortURL, adminURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = a.repo.NewLongURL(a.ctx, a.dbpool, longURL, shortURLID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURLdisplay, adminURLdisplay := fmt.Sprintf("/s/%d/%s", shortURLID, shortURL), fmt.Sprintf("/a/%d/%s", shortURLID, adminURL)

	data.Content = []struct {
		Title string
		Link  string
	}{
		{"Ваша ссылка", longURL},
		{"Короткая ссылка", shortURLdisplay},
		{"Админская ссылка", adminURLdisplay},
	}

	a.IndexPage(w, data)
}

func (a app) LongToShort(w http.ResponseWriter, r *http.Request) {
	shortURLID := chi.URLParam(r, "ID")
	shortURLCode := chi.URLParam(r, "code")
	ip := r.RemoteAddr

	longURL, err := a.repo.GetLongURLByShortIDAndCode(a.ctx, a.dbpool, shortURLID, shortURLCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = a.repo.NewShortURLUsage(a.ctx, a.dbpool, ip, shortURLID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, longURL.LongURL, http.StatusSeeOther)
	return
}

func (a app) AdminsPage(w http.ResponseWriter, r *http.Request) {
	adminURLID := chi.URLParam(r, "ID")
	adminURLCode := chi.URLParam(r, "code")
	data := struct {
		Err   bool
		Link  string
		Count int
	}{}

	admin := filepath.Join("public", "html", "admin.html")
	common := filepath.Join("public", "html", "common.html")
	tmpl, err := template.ParseFiles(admin, common)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := a.repo.GetShortURLByAdminIDAndCode(a.ctx, a.dbpool, adminURLID, adminURLCode)
	if err != nil {
		data.Err = true
		err = tmpl.ExecuteTemplate(w, "admin", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	count, err := a.repo.GetLongURLCountByAdminIDAndCode(a.ctx, a.dbpool, adminURLID, adminURLCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.Link = fmt.Sprintf("/s/%d/%s", shortURL.ID, shortURL.ShortURLCode)
	data.Count = count
	err = tmpl.ExecuteTemplate(w, "admin", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool) *app {
	return &app{ctx, dbpool, repository.NewRepository(dbpool)}
}
