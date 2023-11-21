package models

import (
	"database/sql"
	"net/http"
	"time"
)

type UserUsecases interface {
	Insert(string, string, string) error
	Authenticate(string, string) (int, error)
	Exists(int) (bool, error)
	GetUserId(*http.Request) (int, error)
	GetUserName(id string) (string, error)
	GetUserInfo(email, name string) (int, error)
	AddToken(int, string) error
	GetToken(string) (string, error)
	IsExpired(string) bool
	RemoveToken(string) error
	IsLogged(*http.Request) bool
	GetUserLikes(int) (map[int]*Post, error)
	GetUserPosts(int) (map[int]*Post, error)
}

type UserRepository interface {
	Insert(string, string, string) error
	Authenticate(string, string) (int, error)
	Exists(int) (bool, error)
	GetUserId(*http.Request) (int, error)
	GetUserName(id string) (string, error)
	GetUserInfo(email, name string) (int, error)
	AddToken(int, string) error
	GetToken(string) (string, error)
	IsExpired(string) (*time.Time, error)
	RemoveToken(string) error
	IsLogged()
	GetUserLikes(int) (map[int]*Post, error)
	GetUserPosts(int) (map[int]*Post, error)
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
}

type UserLoginForm struct {
	Email    string
	Password string
}

type UserLikeData struct {
	ID      int  `json:"postID"`
	Likes   int  `json:"likeCount"`
	IsLiked bool `json:"isLiked"`
}

type UserDislikeData struct {
	ID         int  `json:"postID"`
	Dislikes   int  `json:"dislikeCount"`
	IsDisliked bool `json:"isDisliked"`
}
