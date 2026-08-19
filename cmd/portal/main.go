package main

import (
	"log"
	"net/http"
	"os"

	"github.com/portfolio/pf-developer-portal/internal/catalog"
	"github.com/portfolio/pf-developer-portal/internal/httpapi"
)

func main() {
	dir := getenv("PORTAL_SPECS_DIR", "specs")
	addr := getenv("PORTAL_HTTP_ADDR", ":8111")
	cat, err := catalog.LoadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("portal listening on %s (%d specs)", addr, len(cat.APIs))
	if err := http.ListenAndServe(addr, httpapi.New(cat)); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
