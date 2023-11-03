package cookies

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

const (
	cookieName = "session"
)

func SetCookie(w http.ResponseWriter, id int) string {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    GetToken(),
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		MaxAge:   3600,
		// Raw:      strconv.Itoa(id),
	}

	http.SetCookie(w, cookie)
	return cookie.Value
}

func GetCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func DeleteCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
		// Raw:      "",
	}
	http.SetCookie(w, cookie)
}

func GetToken() string {
	token, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return token.String()
}
