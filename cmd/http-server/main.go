package main

import (
	"fmt"
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

func listenAddress() string {
	return fmt.Sprintf("%s:%s", host(), port())
}

func main() {
	address := listenAddress()
	fmt.Printf("HTTP Server Ready at http://%s\n", address)

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(address, nil)
}
