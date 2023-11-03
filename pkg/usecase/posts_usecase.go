package usecase

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"forum.bbilisbe/internal/models"
)

type postsUsecase struct {
	postsRepo models.PostRepository
	usersRepo models.UserRepository
}

func NewPostUsecase(p models.PostRepository, u models.UserRepository) models.PostUsecases {
	return &postsUsecase{
		postsRepo: p,
		usersRepo: u,
	}
}

func (m *postsUsecase) GetPostId(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	id, err := 0, errors.New("")

	if len(parts) >= 4 || parts[1] == "post" {
		id, err = strconv.Atoi(parts[3])
		if err != nil || id < 1 {
			return 0, err
		}
	}
	return id, nil
}

func (m *postsUsecase) Insert(title, content, author string, categories []string) (int, error) {
	return m.postsRepo.Insert(title, content, author, categories)
}

func (m *postsUsecase) CategoryInsert(postid int64, categories []string) error {
	return m.postsRepo.CategoryInsert(postid, categories)
}

// This will insert a new post into the database.
func (m *postsUsecase) Get(id int) (*models.Post, error) {
	return m.postsRepo.Get(id)
}

// This will return the 10 most recently created posts.

func (m *postsUsecase) Latest() (map[int]*models.Post, error) {
	return m.postsRepo.Latest()
}

func (m *postsUsecase) FilteredPosts(categories []string) (map[int]*models.Post, error) {
	posts, err := m.postsRepo.FilteredPosts(categories)
	if err != nil {
		return nil, err
	}

	// Iterate through each post and check if all categories are present in the tags.
	filteredPosts := make(map[int]*models.Post)
	for key, post := range posts {
		postCategories := strings.Split(post.Tags, " ")

		if containsAllCategories(postCategories, categories) {
			filteredPosts[key] = post
		}
	}

	return filteredPosts, nil
}

func containsAllCategories(tags []string, categories []string) bool {
	categorySet := make(map[string]struct{})

	// Create a set of categories for quick look-up.
	for _, category := range tags {
		categorySet[category] = struct{}{}
	}
	// Check if all categories are in the tag set.
	for _, tag := range categories {
		if _, ok := categorySet[tag]; !ok {
			return false
		}
	}

	return true
}

func (m *postsUsecase) LikeInsert(likeData models.UserLikeData, likedBy string) error {
	return m.postsRepo.LikeInsert(likeData, likedBy)
}

func (m *postsUsecase) DislikeInsert(dislikeData models.UserDislikeData, likedBy string) error {
	return m.postsRepo.DislikeInsert(dislikeData, likedBy)
}

func (m *postsUsecase) IsLikedByUser(user string, postid int) bool {
	return m.postsRepo.IsLikedByUser(user, postid)
}

func (m *postsUsecase) IsDislikedByUser(user string, postid int) bool {
	return m.postsRepo.IsDislikedByUser(user, postid)
}

func (m *postsUsecase) CommentInsert(comment, commentBy string, postId int) error {
	return m.postsRepo.CommentInsert(comment, commentBy, postId)
}

func (m *postsUsecase) GetComments(postId int, user string) ([]*models.PostComments, error) {
	return m.postsRepo.GetComments(postId, user)
}

func (m *postsUsecase) CommentLikeInsert(likeData models.CommentLikeData, likedBy string) error {
	return m.postsRepo.CommentLikeInsert(likeData, likedBy)
}

func (m *postsUsecase) IsCommentLikedByUser(user string, commentid int) bool {
	return m.postsRepo.IsCommentLikedByUser(user, commentid)
}

func (m *postsUsecase) CommentDislikeInsert(dislikeData models.CommentDislikeData, dislikedBy string) error {
	return m.postsRepo.CommentDislikeInsert(dislikeData, dislikedBy)
}

func (m *postsUsecase) IsCommentDislikedByUser(user string, commentid int) bool {
	return m.postsRepo.IsCommentDislikedByUser(user, commentid)
}
