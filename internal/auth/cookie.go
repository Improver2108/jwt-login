package auth

import (
	"errors"
	"net/http"
	"time"
)

func SetAuthCookie(w http.ResponseWriter, t *Tokens) {
	cookie := &http.Cookie{
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}

	cookie.Name = "access_token"
	cookie.Value = t.Access
	cookie.MaxAge = int(time.Until(t.ExpAcc).Seconds())
	http.SetCookie(w, cookie)

	cookie.Name = "refresh_token"
	cookie.Value = t.Refresh
	cookie.MaxAge = int(time.Until(t.ExpRef).Seconds())
	http.SetCookie(w, cookie)
}

func ClearAuthCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}

	cookie.Name = "access_token"
	cookie.Value = ""
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	cookie.Name = "refresh_token"
	cookie.Value = ""
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}

func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil || cookie.Value == "" {
		return "", errors.New("missing cookie: " + name)
	}
	return cookie.Value, nil
}
