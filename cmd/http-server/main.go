package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func host() string {
	host := os.Getenv("HOST")
	if host == "" {
		return "localhost"
	}
	return host
}

func port() string {
	host := os.Getenv("PORT")
	if host == "" {
		return "8080"
	}
	return host
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
	http.Handle("/", http.FileServer(http.Dir(root)))

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%s %s\n", r.Method, r.URL)
			http.DefaultServeMux.ServeHTTP(w, r)
		},
	)
}

func startServer(address string, directory string) {
	fmt.Printf("Starting Server...!\n\thttp://%s\n", address)

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
