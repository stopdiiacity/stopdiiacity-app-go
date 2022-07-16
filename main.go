package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/swaggo/http-swagger"

	"github.com/stopdiiacity/stopdiiacity-app-go/api"

	_ "github.com/stopdiiacity/stopdiiacity-app-go/apidocs"
)

//go:embed public
var static embed.FS

func main() {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler())
	r.Post("/verify.json", api.VerifyHandler)
	r.Get("/count.json", api.CountHandler)
	r.Get("/links.json", api.LinksHandler)
	r.Mount("/", http.FileServer(must()))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}

func must() http.FileSystem {
	sub, err := fs.Sub(static, "public")

	if err != nil {
		panic(err)
	}

	return http.FS(sub)
}
