package auth

import (
	"context"
	"net/http"
)

// ログインのresponse.writer
type signInContextKey string

const signInWriterKey signInContextKey = "signInWriter"
const authCookieKey string = "token"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: ログイン時のCookieをResponseWriterでセットするためのcontextをセット
		ctx := context.WithValue(r.Context(), signInWriterKey, w)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func SetAuthCookie(ctx context.Context, token string) {
	w, _ := ctx.Value(signInWriterKey).(http.ResponseWriter)

	week := 60 * 60 * 24 * 7

	cookie := http.Cookie{
		HttpOnly: true,
		MaxAge:   week * 2,
		Secure:   true,
		Name:     authCookieKey,
		Value:    token,
	}
	http.SetCookie(w, &cookie)
}
