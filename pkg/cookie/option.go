package cookie

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/xlgmokha/x/pkg/crypt"
	"github.com/xlgmokha/x/pkg/x"
)

func With(with x.Configure[*http.Cookie]) x.Option[*http.Cookie] {
	return x.With[*http.Cookie](with)
}

func WithValue(value string) x.Option[*http.Cookie] {
	return With(func(c *http.Cookie) {
		c.Value = value
	})
}

func WithSignedValue(value string, signer crypt.Signer) x.Option[*http.Cookie] {
	signature, _ := signer.Sign([]byte(value))
	delimiter := "--"
	return WithValue(fmt.Sprintf("%v%v%v", value, delimiter, string(signature)))
}

func WithEncodedValue(data []byte, encoding *base64.Encoding) x.Option[*http.Cookie] {
	return WithValue(encoding.EncodeToString(data))
}

func WithBase64Value(data []byte) x.Option[*http.Cookie] {
	return WithEncodedValue(data, base64.URLEncoding)
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
