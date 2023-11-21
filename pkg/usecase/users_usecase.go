package usecase

import (
	"net/http"
	"time"

	"forum.bbilisbe/internal/cookies"
	"forum.bbilisbe/internal/models"
)

type userUsecase struct {
	postsRepo models.PostRepository
	usersRepo models.UserRepository
}

func NewUserUsecase(p models.PostRepository, u models.UserRepository) models.UserUsecases {
	return &userUsecase{
		postsRepo: p,
		usersRepo: u,
	}
}

func (m *userUsecase) IsLogged(r *http.Request) bool {
	cookie, errC := cookies.GetCookie(r)

	var data bool
	if errC != nil || cookie.Value == "" {
		data = false
	} else {
		data = true
	}
	return data
}

func (m *userUsecase) Insert(username, email, password string) error {
	return m.usersRepo.Insert(username, email, password)
}

func (m *userUsecase) Authenticate(email, password string) (int, error) {
	return m.usersRepo.Authenticate(email, password)
}

func (m *userUsecase) Exists(id int) (bool, error) {
	return m.usersRepo.Exists(id)
}

func (m *userUsecase) AddToken(id int, token string) error {
	return m.usersRepo.AddToken(id, token)
}

func (m *userUsecase) RemoveToken(token string) error {
	return m.usersRepo.RemoveToken(token)
}

func (m *userUsecase) GetToken(token string) (string, error) {
	return m.usersRepo.GetToken(token)
}

func (m *userUsecase) IsExpired(token string) bool {
	t, _ := m.usersRepo.IsExpired(token)
	return t.Before(time.Now())

}

func (m *userUsecase) GetUserId(r *http.Request) (int, error) {
	return m.usersRepo.GetUserId(r)
}
func (m *userUsecase) GetUserName(id string) (string, error) {
	return m.usersRepo.GetUserName(id)
}

func (m *userUsecase) GetUserInfo(email, name string) (int, error) {
	return m.usersRepo.GetUserInfo(email, name)
 }

func (m *userUsecase) GetUserPosts(author int) (map[int]*models.Post, error) {
	return m.usersRepo.GetUserPosts(author)
}

func (m *userUsecase) GetUserLikes(user int) (map[int]*models.Post, error) {
	return m.usersRepo.GetUserLikes(user)
}
