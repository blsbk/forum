package models

import (
	"database/sql"
	"net/http"
	"time"
)

type PostUsecases interface {
	Insert(PostCreateForm, int) (int, error)
	Get(int) (*Post, error)
	Latest() (map[int]*Post, error)
	GetPostId(*http.Request) (int, error)
	FilteredPosts([]string) (map[int]*Post, error)
	LikeInsert(UserLikeData, int) error
	DislikeInsert(UserDislikeData, int) error
	IsLikedByUser(int, int) bool
	IsDislikedByUser(int, int) bool
	CommentInsert(string, int, int) error
	GetComments(int, int) ([]*PostComments, error)
	CommentLikeInsert(CommentLikeData, int) error
	IsCommentLikedByUser(int, int) bool
	CommentDislikeInsert(CommentDislikeData, int) error
	IsCommentDislikedByUser(int, int) bool
	CategoryInsert(int64, []string) error
}

type PostRepository interface {
	Insert(PostCreateForm, int) (int, error)
	Get(int) (*Post, error)
	Latest() (map[int]*Post, error)
	GetPostId()
	FilteredPosts([]string) (map[int]*Post, error)
	LikeInsert(UserLikeData, int) error
	DislikeInsert(UserDislikeData, int) error
	IsLikedByUser(int, int) bool
	IsDislikedByUser(int, int) bool
	CommentInsert(string, int, int) error
	GetComments(int, int) ([]*PostComments, error)
	CommentLikeInsert(CommentLikeData, int) error
	IsCommentLikedByUser(int, int) bool
	CommentDislikeInsert(CommentDislikeData, int) error
	IsCommentDislikedByUser(int, int) bool
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
	Image    string
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
	ImageURL   string
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
