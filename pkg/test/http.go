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

func Request(method, target string, options ...x.Option[*http.Request]) *http.Request {
	request := httptest.NewRequest(method, target, nil)
	for _, option := range options {
		request = option(request)
	}
	return request
}

func RequestResponse(method, target string, options ...x.Option[*http.Request]) (*http.Request, *httptest.ResponseRecorder) {
	return Request(method, target, options...), httptest.NewRecorder()
}

func WithAcceptHeader(value serde.MediaType) x.Option[*http.Request] {
	return WithRequestHeader("Accept", value.String())
}

func WithRequestHeader(key, value string) x.Option[*http.Request] {
	return func(r *http.Request) *http.Request {
		r.Header.Set(key, value)
		return r
	}
}

func WithContentType[T any](item T, mediaType serde.MediaType) x.Option[*http.Request] {
	body := bytes.NewBuffer(nil)
	x.Check(serde.To[T](body, item, mediaType))
	return WithRequestBody(io.NopCloser(body))
}

func WithRequestBody(body io.ReadCloser) x.Option[*http.Request] {
	return func(r *http.Request) *http.Request {
		r.Body = body
		return r
	}
}

func WithContext(ctx context.Context) x.Option[*http.Request] {
	return func(r *http.Request) *http.Request {
		return r.WithContext(ctx)
	}
}

func WithContextKeyValue[T any](ctx context.Context, key xcontext.Key[T], item T) x.Option[*http.Request] {
	return WithContext(key.With(ctx, item))
}

func WithCookie(cookie *http.Cookie) x.Option[*http.Request] {
	return func(r *http.Request) *http.Request {
		r.AddCookie(cookie)
		return r
	}
}
