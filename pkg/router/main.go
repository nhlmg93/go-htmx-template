package router

import (
	"embed"
	"github.com/nhlmg93/go-htmx-template/pkg/web"
	"net/http"
	"text/template"
)

var (
	Router *http.ServeMux
	html   *template.Template
)

func init() {
	Router = http.NewServeMux()
	Router.Handle("GET /", web.Action(index))
	Router.Handle("GET /index.html", web.Action(index))
}

func SetHtmlTemplates(templateFS *embed.FS) {
	var err error
	html, err = web.TemplateParseFSRecursive(templateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}
}

func index(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, html, "index.html", "", nil)
}
