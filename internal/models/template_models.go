package models

import "forum.bbilisbe/internal/validator"

type TemplateData struct {
	CurrentYear int
	Post        *Post
	Posts       map[int]*Post
	Form        any
	Logged      bool
	IsLiked     bool
	IsDisliked  bool
	Comments    []*PostComments
	validator.Validator
}
