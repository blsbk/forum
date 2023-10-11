package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"forum.bbilisbe/internal/cookies"
	"golang.org/x/crypto/bcrypt"
)

func (m *UserModel) Insert(username, email, password string) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (username, email, hashed_password, created)
	VALUES(?, ?, ?, datetime('now'))`

	_, err = m.DB.Exec(stmt, username, email, string(hashedPwd))
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return ErrDuplicateEmail
		}

		return err
		// var SQLError *sqlite3.Error
		// if errors.As(err, &SQLError) {
		// 	if SQLError.Code == 1062 && strings.Contains(SQLError.err, "users_uc_email") {
		// 	}
		// }
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPwd []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPwd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPwd, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) AddToken(id int, token string) error {
	stmt := `UPDATE users SET token = ?, expiry = DATETIME('now', '+1 hours')
	WHERE ? = id`

	_, err := m.DB.Exec(stmt, token, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) RemoveToken(token string) error {
	stmt := `UPDATE users SET token = NULL, expiry = NULL WHERE token = ?`
	_, err := m.DB.Exec(stmt, token)
	return err
}

func (m *UserModel) GetUserInfo(r *http.Request) (string, error) {
	var author string

	cookie, err := cookies.GetCookie(r)
	if err != nil {
		return "", err
	}
	token := cookie.Value

	stmt := "SELECT username FROM users WHERE token = ?"

	err = m.DB.QueryRow(stmt, token).Scan(&author)
	if err != nil {
		return "", err
	}
	return author, nil
}

func (m *UserModel) IsLogged(r *http.Request) bool {
	_, errC := cookies.GetCookie(r)

	var data bool
	if errC != nil {
		data = false
	} else {
		data = true
	}
	return data
}

func (m *PostModel) GetUserLikes(user string) (map[int]*Post, error) {
	stmt := `SELECT id, title, content, created, author, likes FROM posts JOIN liked ON posts.id = liked.postid WHERE liked.by = ?`

	rows, err := m.DB.Query(stmt, user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer rows.Close()

	posts := map[int]*Post{}

	for rows.Next() {
		p := &Post{}

		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Author, &p.Likes)
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
