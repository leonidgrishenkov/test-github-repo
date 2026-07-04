// Package main is a tiny HTTP server used to exercise the shared docker workflow:
// multi-arch build, SBOM/provenance, cosign sign, and a `--version` smoke test.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	ver := flag.Bool("version", false, "print version and exit")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	if *ver {
		fmt.Println("hello-docker", version)
		return
	}

	msg := os.Getenv("GREETING")
	if msg == "" {
		msg = "Hello from leonidgrishenkov/test-github-repo"
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, msg)
	})

	log.Printf("hello-docker %s listening on %s", version, *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
