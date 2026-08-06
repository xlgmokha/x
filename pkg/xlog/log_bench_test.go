package xlog

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// These two benchmarks exist to gate the zerolog to log/slog port. Both call
// only xlog's own API, so the same source compiles against either backend and
// the numbers are comparable. BenchmarkHTTP covers WithFields transitively.

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		New(io.Discard, Fields{"service": "bench"})
	}
}

func BenchmarkHTTP(b *testing.B) {
	logger := New(io.Discard, Fields{"service": "bench"})
	handler := HTTP(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	request := httptest.NewRequest("GET", "/example", nil)
	recorder := httptest.NewRecorder()

	b.ReportAllocs()
	for b.Loop() {
		handler.ServeHTTP(recorder, request)
	}
}
