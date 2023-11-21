package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"forum.bbilisbe/internal/cookies"
	"forum.bbilisbe/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type sqlUserRepository struct {
	Conn *sql.DB
}

func NewSqlUsersRepository(conn *sql.DB) models.UserRepository {
	return &sqlUserRepository{conn}
}

func (m *sqlUserRepository) Insert(username, email, password string) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (username, email, hashed_password, created)
	VALUES(?, ?, ?, datetime('now'))`

	_, err = m.Conn.Exec(stmt, username, email, string(hashedPwd))
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return models.ErrDuplicateEmail
		} else if err.Error() == "UNIQUE constraint failed: users.username" {
			return models.ErrDuplicateUsername
		}

		return err
	}

	return nil
}

func (m *sqlUserRepository) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPwd []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.Conn.QueryRow(stmt, email).Scan(&id, &hashedPwd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("here1")
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPwd, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			fmt.Println("here2")

			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *sqlUserRepository) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`

	err := m.Conn.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *sqlUserRepository) AddToken(id int, token string) error {
	stmt := `UPDATE users SET token = ?, expiry = DATETIME('now', '+1 hours')
	WHERE ? = id`

	_, err := m.Conn.Exec(stmt, token, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *sqlUserRepository) RemoveToken(token string) error {
	stmt := `UPDATE users SET token = NULL, expiry = NULL WHERE token = ?`
	_, err := m.Conn.Exec(stmt, token)
	return err
}

func (m *sqlUserRepository) GetUserId(r *http.Request) (int, error) {
	var user int

	cookie, err := cookies.GetCookie(r)
	if err != nil {
		return 0, err
	}
	token := cookie.Value

	stmt := "SELECT id FROM users WHERE token = ?"

	err = m.Conn.QueryRow(stmt, token).Scan(&user)
	if err != nil {
		return 0, err
	}
	return user, nil
}

func (m *sqlUserRepository) GetUserName(id string) (string, error) {
	var user string

	stmt := "SELECT username FROM users WHERE id = ?"

	err := m.Conn.QueryRow(stmt, id).Scan(&user)
	if err != nil {
		return "", err
	}
	return user, nil
}

func (m *sqlUserRepository) GetUserInfo(email, name string) (int, error) {
	var user int

	stmt := "SELECT id FROM users WHERE email = ? AND username = ?"

	err := m.Conn.QueryRow(stmt, email, name).Scan(&user)
	if err != nil {
		return 0, err
	}
	return user, nil
}

func (m *sqlUserRepository) GetToken(token string) (string, error) {
	var result bool
	var session string
	stmt := `SELECT EXISTS(SELECT true FROM users WHERE token = ?)`

	err := m.Conn.QueryRow(stmt, token).Scan(&result)
	if err != nil {
		return "", err
	}
	if result {
		stmt2 := `SELECT token FROM users WHERE token = ?`
		err = m.Conn.QueryRow(stmt2, token).Scan(&session)
		if err != nil {
			return "", err
		}
	}
	return session, nil
}

func (m *sqlUserRepository) IsExpired(token string) (*time.Time, error) {
	var result *time.Time
	stmt := `SELECT expiry FROM users WHERE token = ?`

	err := m.Conn.QueryRow(stmt, token).Scan(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *sqlUserRepository) IsLogged() {
	return
}

func (m *sqlUserRepository) GetUserPosts(author int) (map[int]*models.Post, error) {
	stmt := `SELECT id, title, content, created, author, likes, tags FROM posts WHERE author = ?`

	rows, err := m.Conn.Query(stmt, author)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := map[int]*models.Post{}

	for rows.Next() {
		p := &models.Post{}

		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Author, &p.Likes, &p.Tags)
		if err != nil {
			return nil, err
		}

		posts[p.ID] = p
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *sqlUserRepository) GetUserLikes(user int) (map[int]*models.Post, error) {
	stmt := `SELECT id, title, content, created, author, likes, tags FROM posts JOIN likes ON posts.id = likes.postid WHERE likes.likedby = ?`

	rows, err := m.Conn.Query(stmt, user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer rows.Close()

	posts := map[int]*models.Post{}

	for rows.Next() {
		p := &models.Post{}

		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Author, &p.Likes, &p.Tags)
		if err != nil {
			return nil, err
		}

		posts[p.ID] = p
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
