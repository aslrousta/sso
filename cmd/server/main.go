package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aslrousta/rand"
	"github.com/gorilla/mux"
	"github.com/leonelquinteros/gotext"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

type webConfig struct {
	Templates  *template.Template
	Minifier   *minify.M
	AssetsPath string
	Hash       string
}

func main() {
	hash, err := rand.RandomString(8, rand.All)
	if err != nil {
		log.Fatal(err)
	}

	assetsPath := os.Getenv("ASSETS_PATH")
	if assetsPath == "" {
		assetsPath, _ = os.Getwd()
	}

	// Initialize locales
	localePath := filepath.Join(assetsPath, "locale")
	gotext.Configure(localePath, "fa_IR", "default")

	templatesPath := filepath.Join(assetsPath, "templates")
	templates, err := template.New("").Funcs(template.FuncMap{
		"t": gotext.Get,
	}).ParseGlob(templatesPath + "/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	minifier := minify.New()
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/javascript", js.Minify)

	c := &webConfig{
		Minifier:   minifier,
		Templates:  templates,
		AssetsPath: assetsPath,
		Hash:       hash,
	}

	r := mux.NewRouter()

	r.HandleFunc("/oauth", oauthHandler(c)).Methods(http.MethodGet)
	r.HandleFunc("/consent", consentHandler(c)).Methods(http.MethodGet)

	// Static routes
	r.HandleFunc("/css", cssHandler(c)).Methods(http.MethodGet)
	r.HandleFunc("/js", jsHandler(c)).Methods(http.MethodGet)

	http.ListenAndServe(":8080", r)
}

func oauthHandler(c *webConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Templates.ExecuteTemplate(w, "login", map[string]string{
			"Hash": c.Hash,
		})
	}
}

func consentHandler(c *webConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Templates.ExecuteTemplate(w, "consent", map[string]string{
			"Hash": c.Hash,
		})
	}
}

func cssHandler(c *webConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join(c.AssetsPath, "css/main.css")
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}
		bytes, err = c.Minifier.Bytes("text/css", bytes)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Add("Content-Type", "text/css")
		w.Header().Add("Cache-Control", "max-age=86400")
		w.Write(bytes)
	}
}

func jsHandler(c *webConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join(c.AssetsPath, "js/main.js")
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}
		bytes, err = c.Minifier.Bytes("text/javascript", bytes)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Add("Content-Type", "text/javascript")
		w.Header().Add("Cache-Control", "max-age=86400")
		w.Write(bytes)
	}
}
