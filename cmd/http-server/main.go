package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/xlgmokha/x/pkg/env"
)

func host() string {
	return env.Fetch("HOST", "localhost")
}

func port() string {
	return env.Fetch("PORT", "8080")
}

func directory() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return "."
}

func listenAddress() string {
	return fmt.Sprintf("%s:%s", host(), port())
}

func buildHttpHandlerFor(root string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(root)))

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%s %s\n", r.Method, r.URL)
			mux.ServeHTTP(w, r)
		},
	)
}

func startServer(address string, directory string) {
	fmt.Printf("Listening and serving HTTP on http://%s\n", address)

	log.Fatal(
		http.ListenAndServe(
			address,
			buildHttpHandlerFor(directory),
		),
	)
}

func main() {
	startServer(listenAddress(), directory())
}
