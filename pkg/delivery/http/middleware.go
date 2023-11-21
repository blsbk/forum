package delivery

import (
	"fmt"
	"net/http"

	"forum.bbilisbe/internal/cookies"
)

func (h *Handler) SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := cookies.GetCookie(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		// Check if the session exists and is valid.
		session, err := h.UUsecase.GetToken(token.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if session == "" || h.UUsecase.IsExpired(session) {
			cookies.DeleteCookie(w)
			h.UUsecase.RemoveToken(session)
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}

		// Continue to the next handler if the session cookie is found and valid.
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				h.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RestrictGet(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			h.clientError(w, http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) RestrictPost(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			h.clientError(w, http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) RestrictGetPost(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodGet)
			w.Header().Add("Allow", http.MethodPost)
			h.clientError(w, http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) RequireLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !h.UUsecase.IsLogged(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	}
}
