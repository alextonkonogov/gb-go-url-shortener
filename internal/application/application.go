package application

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/dchest/uniuri"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"

	"github.com/alextonkonogov/gb-go-url-shortener/internal/repository"
)

const URLregexp = `(?m)https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&\/\/=]*)`

type app struct {
	ctx    context.Context
	dbpool *pgxpool.Pool
	repo   *repository.Repository
}

func (a app) Routes(r *httprouter.Router) {
	r.ServeFiles("/public/*filepath", http.Dir("public"))
	r.GET("/", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		a.IndexPage(rw, nil)
	})
	r.POST("/short", a.ShortURL)
	r.GET("/s/:ID/:code", a.LongToShort)
	r.GET("/a/:ID/:code", a.AdminsPage)
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

func (a app) ShortURL(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
		a.IndexPage(rw, data)
		return
	}

	shortURL, adminURL = uniuri.New(), uniuri.New()

	shortURLID, err := a.repo.NewShortURL(a.ctx, a.dbpool, shortURL, adminURL)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = a.repo.NewLongURL(a.ctx, a.dbpool, longURL, shortURLID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
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

	a.IndexPage(rw, data)
}

func (a app) LongToShort(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	shortURLID := p.ByName("ID")
	shortURLCode := p.ByName("code")
	ip := r.Header.Get("X-FORWARDED-FOR")

	longURL, err := a.repo.GetLongURLByShortIDAndCode(a.ctx, a.dbpool, shortURLID, shortURLCode)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = a.repo.NewShortURLUsage(a.ctx, a.dbpool, ip, shortURLID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(rw, r, longURL.LongURL, http.StatusSeeOther)
	return
}

func (a app) AdminsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	adminURLID := p.ByName("ID")
	adminURLCode := p.ByName("code")
	data := struct {
		Err   bool
		Link  string
		Count int
	}{}

	admin := filepath.Join("public", "html", "admin.html")
	common := filepath.Join("public", "html", "common.html")
	tmpl, err := template.ParseFiles(admin, common)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := a.repo.GetShortURLByAdminIDAndCode(a.ctx, a.dbpool, adminURLID, adminURLCode)
	if err != nil {
		data.Err = true
		err = tmpl.ExecuteTemplate(rw, "admin", data)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	count, err := a.repo.GetLongURLCountByAdminIDAndCode(a.ctx, a.dbpool, adminURLID, adminURLCode)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	data.Link = fmt.Sprintf("/s/%d/%s", shortURL.ID, shortURL.ShortURLCode)
	data.Count = count
	err = tmpl.ExecuteTemplate(rw, "admin", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool) *app {
	return &app{ctx, dbpool, repository.NewRepository(dbpool)}
}
