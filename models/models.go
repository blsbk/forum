package models

import (
	"database/sql"
	"time"

	"forum.bbilisbe/internal/validator"
)

type Post struct {
	ID      int `json:"postID"`
	Title   string
	Content string
	Created time.Time
	Author  string
	Likes   int `json:"likeCount"`
}

type PostComments struct {
	Id      int
	Comment string
	Author  string
	Likes   int
	IsLiked bool
}

type PostModel struct {
	DB *sql.DB
}

type PostCreateForm struct {
	Title      string
	Content    string
	Categories []string
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Token          *string
	Expiry         *time.Time
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

type UserSignupForm struct {
	Name     string
	Email    string
	Password string
	validator.Validator
}

type UserLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

type UserLikeData struct {
	ID      int  `json:"postID"`
	Likes   int  `json:"likeCount"`
	IsLiked bool `json:"isLiked"`
}

type CommentLikeData struct {
	ID      int  `json:"commentID"`
	PostId  int  `json:"postID"`
	Likes   int  `json:"commentLikeCount"`
	IsLiked bool `json:"isCommentLiked"`
}
