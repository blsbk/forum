package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"forum.bbilisbe/internal/cookies"
	"forum.bbilisbe/internal/validator"
	"forum.bbilisbe/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	posts, err := app.posts.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		categories := r.Form["category"]

		posts, err = app.posts.FilteredPosts(categories)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	data := app.newTemplateData(r)
	data.Posts = posts

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	postId, err := app.posts.GetPostId(r)
	if err != nil {
		app.notFound(w)
		return
	}

	post, err := app.posts.Get(postId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	Comments, _ := app.posts.GetComments(postId)
	data.Comments = Comments
	data.Post = post
	user, err := app.users.GetUserInfo(r)

	if err == nil {
		data.IsLiked = app.posts.IsLikedByUser(user, postId)
		if len(Comments) != 0 {
			for _, comm := range data.Comments {
				comm.IsLiked = app.posts.IsCommentLikedByUser(user, comm.Id)
			}
		}
	}

	if r.Method == http.MethodPost {
		comment := models.PostComments{
			Comment: r.FormValue("comment"),
		}
		data.CheckField(validator.NotBlank(comment.Comment), "comment", "This field cannot be blank")
		data.CheckField(validator.MaxChars(comment.Comment, 100), "comment", "This field cannot be more than 100 characters long")
		if !data.Valid() {
			app.render(w, http.StatusUnprocessableEntity, "view.html", data)
			return
		}

		err := app.posts.CommentInsert(comment.Comment, user, postId)
		if err != nil {
			app.serverError(w, err)
			return
		}

		Comments, _ := app.posts.GetComments(postId)
		fmt.Println(Comments)
		data.Comments = Comments
	}

	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) postCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		data := app.newTemplateData(r)
		data.Form = models.PostCreateForm{}
		app.render(w, http.StatusOK, "create.html", data)

	} else if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		data := app.newTemplateData(r)

		form := models.PostCreateForm{
			Title:      r.PostForm.Get("title"),
			Content:    r.PostForm.Get("content"),
			Categories: r.Form["category"],
		}

		data.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		data.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		data.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

		if !data.Valid() {
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "create.html", data)
			return
		}

		author, _ := app.users.GetUserInfo(r)
		id, err := app.posts.Insert(form.Title, form.Content, author, form.Categories)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
	}
}

func (app *application) postLike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/post/like" {

		user, _ := app.users.GetUserInfo(r)

		var likeData models.UserLikeData

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&likeData); err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		err := app.posts.LikeInsert(likeData, user)
		if err != nil {
			app.serverError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (app *application) commentLike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/post/commentLike" {

		user, _ := app.users.GetUserInfo(r)

		var commentLikeData models.CommentLikeData

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&commentLikeData); err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		err := app.posts.CommentLikeInsert(commentLikeData, user)
		if err != nil {
			app.serverError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := app.newTemplateData(r)
		data.Form = models.UserSignupForm{}
		app.render(w, http.StatusOK, "signup.html", data)
	} else if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		data := app.newTemplateData(r)

		form := models.UserSignupForm{
			Name:     r.PostForm.Get("name"),
			Email:    r.PostForm.Get("email"),
			Password: r.PostForm.Get("password"),
		}

		data.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
		data.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		data.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
		data.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
		data.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

		if !data.Valid() {
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
			return
		}

		err = app.users.Insert(form.Name, form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				form.AddFieldError("email", "Email address is already in use")

				data := app.newTemplateData(r)
				data.Form = form
				app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
			} else {
				app.serverError(w, err)
			}
			return
		}

		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := app.newTemplateData(r)
		data.Form = models.UserLoginForm{}
		app.render(w, http.StatusOK, "login.html", data)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		data := app.newTemplateData(r)

		form := models.UserLoginForm{
			Email:    r.PostForm.Get("email"),
			Password: r.PostForm.Get("password"),
		}

		data.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		data.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
		data.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

		if !data.Valid() {
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
			return
		}

		id, err := app.users.Authenticate(form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				form.AddNonFieldError("Email or password is incorrect")
				data := app.newTemplateData(r)
				data.Form = form
				app.render(w, http.StatusUnprocessableEntity, "login.html", data)
			} else {
				app.serverError(w, err)
			}
			return
		}

		token := cookies.SetCookie(w)
		app.users.AddToken(id, token)

		http.Redirect(w, r, "/post/create", http.StatusSeeOther)
	}
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := cookies.GetCookie(r)
	if cookie != nil {
		if err != nil {
			app.serverError(w, err)
		}
		app.users.RemoveToken(cookie.Value)
		cookies.DeleteCookie(w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *application) userPosts(w http.ResponseWriter, r *http.Request) {
	author, _ := app.users.GetUserInfo(r)

	posts, err := app.posts.GetUserPosts(author)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Posts = posts

	app.render(w, http.StatusOK, "userposts.html", data)
}

func (app *application) userLikedPosts(w http.ResponseWriter, r *http.Request) {
	user, _ := app.users.GetUserInfo(r)

	posts, err := app.posts.GetUserLikes(user)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Posts = posts

	app.render(w, http.StatusOK, "userposts.html", data)
}
