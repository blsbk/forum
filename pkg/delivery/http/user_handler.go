package delivery

import (
	"errors"
	"net/http"

	"forum.bbilisbe/internal/cookies"
	"forum.bbilisbe/internal/models"
	"forum.bbilisbe/internal/validator"
)

func (h *Handler) userSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := h.newTemplateData(r)
		if data.Logged {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		data.Form = models.UserSignupForm{}
		h.render(w, http.StatusOK, "signup.html", data)
	} else if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		data := h.newTemplateData(r)

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
			h.render(w, http.StatusUnprocessableEntity, "signup.html", data)
			return
		}

		err = h.UUsecase.Insert(form.Name, form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {

				data := h.newTemplateData(r)
				data.AddFieldError("email", "Email address is already in use")
				data.Form = form
				h.render(w, http.StatusUnprocessableEntity, "signup.html", data)
			} else if errors.Is(err, models.ErrDuplicateUsername) {

				data := h.newTemplateData(r)
				data.AddFieldError("name", "Username is already in use")
				data.Form = form
				h.render(w, http.StatusUnprocessableEntity, "signup.html", data)
			} else {
				h.serverError(w, err)
			}
			return
		}

		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

func (h *Handler) userLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := h.newTemplateData(r)
		if data.Logged {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		data.Form = models.UserLoginForm{}
		h.render(w, http.StatusOK, "login.html", data)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		data := h.newTemplateData(r)

		form := models.UserLoginForm{
			Email:    r.PostForm.Get("email"),
			Password: r.PostForm.Get("password"),
		}

		data.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		data.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
		data.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

		if !data.Valid() {
			data.Form = form
			h.render(w, http.StatusUnprocessableEntity, "login.html", data)
			return
		}

		id, err := h.UUsecase.Authenticate(form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				data := h.newTemplateData(r)
				data.AddNonFieldError("email", "Email or password is incorrect")
				data.Form = form
				h.render(w, http.StatusUnprocessableEntity, "login.html", data)
			} else {
				h.serverError(w, err)
			}
			return
		}

		token := cookies.SetCookie(w, id)
		h.UUsecase.AddToken(id, token)

		http.Redirect(w, r, "/post/create", http.StatusSeeOther)
	}
}

func (h *Handler) userLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := cookies.GetCookie(r)
	if cookie != nil {
		if err != nil {
			h.serverError(w, err)
		}
		h.UUsecase.RemoveToken(cookie.Value)
		cookies.DeleteCookie(w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) userPosts(w http.ResponseWriter, r *http.Request) {
	author, _ := h.UUsecase.GetUserInfo(r)

	posts, err := h.UUsecase.GetUserPosts(author)
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := h.newTemplateData(r)
	data.Posts = posts

	h.render(w, http.StatusOK, "userposts.html", data)
}

func (h *Handler) userLikedPosts(w http.ResponseWriter, r *http.Request) {
	user, _ := h.UUsecase.GetUserInfo(r)

	posts, err := h.UUsecase.GetUserLikes(user)
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := h.newTemplateData(r)
	data.Posts = posts

	h.render(w, http.StatusOK, "userposts.html", data)
}
