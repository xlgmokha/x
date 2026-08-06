package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/xlgmokha/x/pkg/env"
	"github.com/xlgmokha/x/pkg/x"
	"github.com/xlgmokha/x/pkg/xlog"
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
	return net.JoinHostPort(host(), port())
}

func buildHttpHandlerFor(root string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(root)))

	return mux
}

func startServer(address string, directory string) {
	logger := xlog.New(os.Stdout, xlog.Fields{})
	handler := x.Middleware[http.Handler](
		buildHttpHandlerFor(directory),
		xlog.HTTP(logger),
	)

	logger.Info("listening", slog.String("address", address), slog.String("directory", directory))

	if err := http.ListenAndServe(address, handler); err != nil {
		logger.Error("server stopped", slog.Any("error", err))
		os.Exit(1)
	}
}

func main() {
	startServer(listenAddress(), directory())
}
