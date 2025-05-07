package cookie

import (
	"net/http"

	"github.com/xlgmokha/x/pkg/x"
)

func Write(w http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(w, cookie)
}

func WriteWith(w http.ResponseWriter, options ...x.Option[*http.Cookie]) {
	Write(w, x.New[*http.Cookie](options...))
}

func Expire(w http.ResponseWriter, name string, options ...x.Option[*http.Cookie]) {
	Write(w, Reset(name, options...))
}
