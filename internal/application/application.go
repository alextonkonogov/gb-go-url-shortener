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
	r.POST("/short-url", a.ShortUrl)
	r.GET("/s/:id/:code", a.LongToShort)
	r.GET("/a/:id/:code", a.AdminsPage)

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

func (a app) ShortUrl(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var (
		longUrl, shortUrl, adminUrl string
	)
	longUrl = r.FormValue("longUrl")

	data := struct {
		Err     bool
		Content template.HTML
	}{}

	if !regexp.MustCompile(URLregexp).MatchString(longUrl) {
		data.Err, data.Content = true, "<p>Вы ввели невалидную ссылку!</p>"
		a.IndexPage(rw, data)
		return
	}

	shortUrl, adminUrl = uniuri.New(), uniuri.New()

	shortUrlId, err := a.repo.NewShortUrl(a.ctx, a.dbpool, shortUrl, adminUrl)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = a.repo.NewLongUrl(a.ctx, a.dbpool, longUrl, shortUrlId)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	shortUrldisplay, adminUrldisplay := fmt.Sprintf("/s/%d/%s", shortUrlId, shortUrl), fmt.Sprintf("/a/%d/%s", shortUrlId, adminUrl)

	data.Content = template.HTML(
		fmt.Sprintf(`<div>
								<h6>Ваша ссылка</h6>
								<a href="%s" target="_blank">%s</a>
							</div>
							<div class="mt-3">
								<h6>Короткая ссылка</h6>
								<a href="%s" target="_blank">%s</a>
							</div>
							<div class="mt-3">
								<h6>Админская ссылка</h6>
								<a href="%s" target="_blank">%s</a>
							</div>`,
			longUrl, longUrl, shortUrldisplay, shortUrldisplay, adminUrldisplay, adminUrldisplay,
		))
	a.IndexPage(rw, data)
}

func (a app) LongToShort(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	shortUrlId := p.ByName("id")
	shortUrlCode := p.ByName("code")
	ip := r.Header.Get("X-FORWARDED-FOR")

	longUrl, err := a.repo.GetLongUrlByShortIdAndCode(a.ctx, a.dbpool, shortUrlId, shortUrlCode)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = a.repo.NewShortUrlUsage(a.ctx, a.dbpool, ip, shortUrlId)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(rw, r, longUrl.LongUrl, http.StatusSeeOther)
	return
}

func (a app) AdminsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	adminUrlId := p.ByName("id")
	adminUrlCode := p.ByName("code")

	shortUrl, err := a.repo.GetShortUrlByAdminIdAndCode(a.ctx, a.dbpool, adminUrlId, adminUrlCode)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	count, err := a.repo.GetLongUrlCountByAdminIdAndCode(a.ctx, a.dbpool, adminUrlId, adminUrlCode)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	data := struct {
		Err     bool
		Content template.HTML
	}{}

	admin := filepath.Join("public", "html", "admin.html")
	common := filepath.Join("public", "html", "common.html")
	tmpl, err := template.ParseFiles(admin, common)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	data.Content = template.HTML(
		fmt.Sprintf(`<div>
								<h6>Ваша ссылка</h6>
								<a href="/s/%d/%s" target="_blank">/s/%d/%s</a>
							</div>
								<div class="mt-3">
									<h6>Кол-во переходов:</h6>
								<p> %d</p>
							</div>`,
			shortUrl.Id, shortUrl.ShortUrlCode, shortUrl.Id, shortUrl.ShortUrlCode, count))
	err = tmpl.ExecuteTemplate(rw, "admin", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool) *app {
	return &app{ctx, dbpool, repository.NewRepository(dbpool)}
}
