package test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"

	xcontext "github.com/xlgmokha/x/pkg/context"
	"github.com/xlgmokha/x/pkg/serde"
	"github.com/xlgmokha/x/pkg/x"
)

type RequestOption x.Option[*http.Request]

func Request(method, target string, options ...RequestOption) *http.Request {
	request := httptest.NewRequest(method, target, nil)
	for _, option := range options {
		request = option(request)
	}
	return request
}

func RequestResponse(method, target string, options ...RequestOption) (*http.Request, *httptest.ResponseRecorder) {
	return Request(method, target, options...), httptest.NewRecorder()
}

func WithAcceptHeader(value serde.MediaType) RequestOption {
	return WithRequestHeader("Accept", string(value))
}

func WithRequestHeader(key, value string) RequestOption {
	return func(r *http.Request) *http.Request {
		r.Header.Set(key, value)
		return r
	}
}

func WithContentType[T any](item T, mediaType serde.MediaType) RequestOption {
	body := bytes.NewBuffer(nil)
	x.Check(serde.To[T](body, item, mediaType))
	return WithRequestBody(io.NopCloser(body))
}

func WithRequestBody(body io.ReadCloser) RequestOption {
	return func(r *http.Request) *http.Request {
		r.Body = body
		return r
	}
}

func WithContext(ctx context.Context) RequestOption {
	return func(r *http.Request) *http.Request {
		return r.WithContext(ctx)
	}
}

func WithContextKeyValue[T any](ctx context.Context, key xcontext.Key[T], item T) RequestOption {
	return WithContext(key.With(ctx, item))
}

func WithCookie(cookie *http.Cookie) RequestOption {
	return func(r *http.Request) *http.Request {
		r.AddCookie(cookie)
		return r
	}
}
