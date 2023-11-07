package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"forum.bbilisbe/internal/models"
	"forum.bbilisbe/internal/validator"
	"github.com/gofrs/uuid"
)

func (h *Handler) postView(w http.ResponseWriter, r *http.Request) {
	postId, err := h.PUsecase.GetPostId(r)
	if err != nil {
		h.notFound(w)
		return
	}

	post, err := h.PUsecase.Get(postId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.notFound(w)
		} else {
			h.serverError(w, err)
		}
		return
	}

	user, err := h.UUsecase.GetUserInfo(r)
	data := h.newTemplateData(r)
	Comments, _ := h.PUsecase.GetComments(postId, user)
	data.Comments = Comments
	data.Post = post

	if err == nil {
		data.IsLiked = h.PUsecase.IsLikedByUser(user, postId)
		data.IsDisliked = h.PUsecase.IsDislikedByUser(user, postId)
		if len(Comments) != 0 {
			for _, comm := range data.Comments {
				comm.IsLiked = h.PUsecase.IsCommentLikedByUser(user, comm.Id)
				comm.IsDisliked = h.PUsecase.IsCommentDislikedByUser(user, comm.Id)
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
			h.render(w, http.StatusUnprocessableEntity, "view.html", data)
			return
		}

		err := h.PUsecase.CommentInsert(comment.Comment, user, postId)
		if err != nil {
			h.serverError(w, err)
			return
		}

		Comments, _ := h.PUsecase.GetComments(postId, user)
		data.Comments = Comments
		http.Redirect(w, r, fmt.Sprintf("/post/view/%d", postId), http.StatusSeeOther)
	}

	h.render(w, http.StatusOK, "view.html", data)
}

func (h *Handler) postCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		data := h.newTemplateData(r)
		data.Form = models.PostCreateForm{}
		h.render(w, http.StatusOK, "create.html", data)

	} else if r.Method == http.MethodPost {

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			fmt.Println(err)
			h.clientError(w, http.StatusBadRequest)
			return
		}
		data := h.newTemplateData(r)

		form := models.PostCreateForm{
			Title:      r.PostForm.Get("title"),
			Content:    r.PostForm.Get("content"),
			Categories: r.Form["category"],
		}

		file, handler, err := r.FormFile("image")
		if err == nil {
			uniqueID, err := uuid.NewV4()
			if err != nil {
				h.serverError(w, err)
				return
			}
			filename := strings.Replace(uniqueID.String(), "-", "", -1)
			prevname := filepath.Base(handler.Filename)
			fileExt := filepath.Ext(prevname)

			if fileExt != ".jpeg" && fileExt != ".png" && fileExt != ".gif" && fileExt != ".jpg" {
				h.clientError(w, http.StatusBadRequest)
			}

			image := fmt.Sprintf("%s%s", filename, fileExt)
			f, err := os.Create(fmt.Sprintf("./ui/static/img/user_images/%s", image))
			if err != nil {
				h.serverError(w, err)
				return
			}
			defer f.Close()

			_, err = io.Copy(f, file)
			if err != nil {
				h.serverError(w, err)
				return
			}
			form.ImageURL = fmt.Sprintf("/static/img/user_images/%s", image)
		}

		if len(form.Categories) == 0 {
			data.Form = form
			data.AddNonFieldError("tags", "Tags cannot be empty")
			h.render(w, http.StatusUnprocessableEntity, "create.html", data)
			return
		}

		data.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		data.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		data.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

		if !data.Valid() {
			data.Form = form
			h.render(w, http.StatusUnprocessableEntity, "create.html", data)
			return
		}

		author, _ := h.UUsecase.GetUserInfo(r)
		id, err := h.PUsecase.Insert(form, author)
		if err != nil {
			h.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
	}
}

func (h *Handler) postLike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/post/like" {

		user, _ := h.UUsecase.GetUserInfo(r)

		var likeData models.UserLikeData

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&likeData); err != nil {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		err := h.PUsecase.LikeInsert(likeData, user)
		if err != nil {
			h.serverError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) postDislike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/post/dislike" {

		user, _ := h.UUsecase.GetUserInfo(r)

		var dislikeData models.UserDislikeData

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&dislikeData); err != nil {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		err := h.PUsecase.DislikeInsert(dislikeData, user)
		if err != nil {
			h.serverError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) commentLike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/post/commentLike" {

		user, _ := h.UUsecase.GetUserInfo(r)

		var commentLikeData models.CommentLikeData

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&commentLikeData); err != nil {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		err := h.PUsecase.CommentLikeInsert(commentLikeData, user)
		if err != nil {
			h.serverError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) commentDislike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/post/commentDislike" {

		user, _ := h.UUsecase.GetUserInfo(r)

		var commentDislikeData models.CommentDislikeData

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&commentDislikeData); err != nil {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		err := h.PUsecase.CommentDislikeInsert(commentDislikeData, user)
		if err != nil {
			h.serverError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
