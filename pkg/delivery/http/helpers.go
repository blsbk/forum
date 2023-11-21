package delivery

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"time"

	"forum.bbilisbe/internal/models"
)

func (h *Handler) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.errorLog.Output(2, trace)
	Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (h *Handler) clientError(w http.ResponseWriter, status int) {
	Errors(w, status, http.StatusText(status))
}

func (h *Handler) notFound(w http.ResponseWriter) {
	h.clientError(w, http.StatusNotFound)
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

func (h *Handler) newTemplateData(r *http.Request) *models.TemplateData {
	return &models.TemplateData{
		CurrentYear: time.Now().Year(),
		Logged:      h.UUsecase.IsLogged(r),
	}
}

func Errors(w http.ResponseWriter, status int, message string) {

	t, err := template.ParseFiles("ui/html/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := struct {
		StatusCodeAndText string
		MessageError      string
	}{
		StatusCodeAndText: strconv.Itoa(status) + " " + http.StatusText(status),
		MessageError:      message,
	}
	w.WriteHeader(status)
	if err := t.Execute(w, res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func generateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}