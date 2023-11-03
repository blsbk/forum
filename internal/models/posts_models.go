package models

import (
	"database/sql"
	"net/http"
	"time"
)

type PostUsecases interface {
	Insert(string, string, string, []string) (int, error)
	Get(int) (*Post, error)
	Latest() (map[int]*Post, error)
	GetPostId(*http.Request) (int, error)
	FilteredPosts([]string) (map[int]*Post, error)
	LikeInsert(UserLikeData, string) error
	DislikeInsert(UserDislikeData, string) error
	IsLikedByUser(string, int) bool
	IsDislikedByUser(string, int) bool
	CommentInsert(string, string, int) error
	GetComments(int, string) ([]*PostComments, error)
	CommentLikeInsert(CommentLikeData, string) error
	IsCommentLikedByUser(string, int) bool
	CommentDislikeInsert(CommentDislikeData, string) error
	IsCommentDislikedByUser(string, int) bool
	CategoryInsert(int64, []string) error
}

type PostRepository interface {
	Insert(string, string, string, []string) (int, error)
	Get(int) (*Post, error)
	Latest() (map[int]*Post, error)
	GetPostId()
	FilteredPosts([]string) (map[int]*Post, error)
	LikeInsert(UserLikeData, string) error
	DislikeInsert(UserDislikeData, string) error
	IsLikedByUser(string, int) bool
	IsDislikedByUser(string, int) bool
	CommentInsert(string, string, int) error
	GetComments(int, string) ([]*PostComments, error)
	CommentLikeInsert(CommentLikeData, string) error
	IsCommentLikedByUser(string, int) bool
	CommentDislikeInsert(CommentDislikeData, string) error
	IsCommentDislikedByUser(string, int) bool
	CategoryInsert(int64, []string) error
}

type Post struct {
	ID       int `json:"postID"`
	Title    string
	Content  string
	Created  time.Time
	Author   string
	Likes    int `json:"likeCount"`
	Dislikes int `json:"dislikeCount"`
	Tags     string
}

type PostComments struct {
	Id         int
	Comment    string
	Author     string
	Likes      int
	Dislikes   int
	IsLiked    bool
	IsDisliked bool
}

type PostModel struct {
	DB *sql.DB
}

type PostCreateForm struct {
	Title      string
	Content    string
	Categories []string
}

type CommentLikeData struct {
	ID      int  `json:"commentID"`
	PostId  int  `json:"postID"`
	Likes   int  `json:"commentLikeCount"`
	IsLiked bool `json:"isCommentLiked"`
}

type CommentDislikeData struct {
	ID         int  `json:"commentID"`
	PostId     int  `json:"postID"`
	Dislikes   int  `json:"commentDislikeCount"`
	IsDisliked bool `json:"isCommentDisliked"`
}
