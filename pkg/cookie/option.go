package cookie

import (
	"net/http"
	"time"

	"github.com/xlgmokha/x/pkg/x"
)

func With(with func(*http.Cookie)) x.Option[*http.Cookie] {
	return func(c *http.Cookie) *http.Cookie {
		with(c)
		return c
	}
}

func WithValue(value string) x.Option[*http.Cookie] {
	return With(func(c *http.Cookie) {
		c.Value = value
	})
}

func WithPath(value string) x.Option[*http.Cookie] {
	return With(func(c *http.Cookie) {
		c.Path = value
	})
}

func WithHttpOnly(value bool) x.Option[*http.Cookie] {
	return With(func(c *http.Cookie) {
		c.HttpOnly = value
	})
}

func WithSecure(value bool) x.Option[*http.Cookie] {
	return With(func(c *http.Cookie) {
		c.Secure = value
	})
}

func WithDomain(value string) x.Option[*http.Cookie] {
	return With(func(c *http.Cookie) {
		c.Domain = value
	})
}

func WithSameSite(value http.SameSite) x.Option[*http.Cookie] {
	return With(func(c *http.Cookie) {
		c.SameSite = value
	})
}

func WithExpiration(expires time.Time) x.Option[*http.Cookie] {
	return With(func(c *http.Cookie) {
		c.Expires = expires
		if expires.Before(time.Now()) {
			c.MaxAge = -1
		} else {
			duration := time.Until(expires).Round(time.Second)
			c.MaxAge = int(duration.Seconds())
		}
	})
}
