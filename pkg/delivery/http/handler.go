package delivery

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"forum.bbilisbe/internal/models"
)

type Handler struct {
	PUsecase      models.PostUsecases
	UUsecase      models.UserUsecases
	templateCache map[string]*template.Template
	infoLog       *log.Logger
	errorLog      *log.Logger
}

func NewPostHandler(pu models.PostUsecases, uu models.UserUsecases, infoLog, errorLog *log.Logger) http.Handler {
	templateCache, _ := newTemplateCache()

	handler := &Handler{
		PUsecase:      pu,
		UUsecase:      uu,
		templateCache: templateCache,
		infoLog:       infoLog,
		errorLog:      errorLog,
	}
	mux := http.NewServeMux()

	// if err != nil {
	// 	errorLog.Fatal(err)
	// }

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", handler.home)
	mux.HandleFunc("/post/view/", handler.postView)
	mux.HandleFunc("/post/create", handler.RequireLog(handler.postCreate))
	mux.HandleFunc("/user/signup", handler.userSignup)
	mux.HandleFunc("/user/login", handler.userLogin)
	mux.HandleFunc("/user/logout", handler.RestrictPost(handler.userLogout))
	mux.HandleFunc("/user/posts", handler.RestrictGet(handler.userPosts))
	mux.HandleFunc("/user/likedposts", handler.RestrictGet(handler.userLikedPosts))
	mux.HandleFunc("/post/like", handler.RequireLog(handler.postLike))
	mux.HandleFunc("/post/dislike", handler.RequireLog(handler.postDislike))
	mux.HandleFunc("/post/commentLike", handler.RequireLog(handler.commentLike))
	mux.HandleFunc("/post/commentDislike", handler.RequireLog(handler.commentDislike))

	return handler.RecoverPanic(handler.AuthMiddleware(handler.LogRequest(handler.SecureHeaders(mux))))
}

func (h *Handler) render(w http.ResponseWriter, status int, page string, data *models.TemplateData) {
	ts, ok := h.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		h.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		h.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.notFound(w)
		return
	}

	posts, err := h.PUsecase.Latest()
	if err != nil {
		h.serverError(w, err)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.clientError(w, http.StatusBadRequest)
			return
		}
		categories := r.Form["category"]

		posts, err = h.PUsecase.FilteredPosts(categories)
		if err != nil {
			h.serverError(w, err)
			return
		}
	}

	data := h.newTemplateData(r)
	data.Posts = posts

	h.render(w, http.StatusOK, "home.html", data)
}
