package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/post/view/", app.postView)
	mux.HandleFunc("/post/create", app.requireLog(app.postCreate))
	mux.HandleFunc("/user/signup", app.userSignup)
	mux.HandleFunc("/user/login", app.userLogin)
	mux.HandleFunc("/user/logout", app.restrictPost(app.userLogout))
	mux.HandleFunc("/user/posts", app.restrictGet(app.userPosts))
	mux.HandleFunc("/user/likedposts", app.restrictGet(app.userLikedPosts))
	mux.HandleFunc("/post/like", app.postLike)
	mux.HandleFunc("/post/commentLike", app.commentLike)
	// mux.HandleFunc("/post/comment", app.postComment)

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
