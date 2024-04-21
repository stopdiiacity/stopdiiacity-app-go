package main

import (
	"crypto/tls"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"

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

	run(r)
}

func must() http.FileSystem {
	sub, err := fs.Sub(static, "public")

	if err != nil {
		panic(err)
	}

	return http.FS(sub)
}

func run(handler http.Handler) {
	runProduction(handler)
}

// https://stackoverflow.com/questions/37321760/how-to-set-up-lets-encrypt-for-a-go-server-application
// https://stackoverflow.com/a/40494806/17655004
func runProduction(handler http.Handler) {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(strings.Split(os.Getenv("HOSTS"), ",")...),
		Cache:      autocert.DirCache(os.Getenv("TLS_CERTIFICATES_DIR")),
	}

	server := &http.Server{
		Addr:    ":https",
		Handler: handler,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
			MinVersion:     tls.VersionTLS12, // improves cert reputation score at https://www.ssllabs.com/ssltest/
		},
	}

	var g errgroup.Group

	g.Go(func() error {
		return http.ListenAndServe(":http", certManager.HTTPHandler(nil))
	})

	g.Go(func() error {
		return server.ListenAndServeTLS("", "") // Key and cert are coming from Let's Encrypt
	})

	log.Fatal(g.Wait())
}

func runDevelopment(handler http.Handler) {
	log.Fatal(http.ListenAndServe(":http", handler))
}
