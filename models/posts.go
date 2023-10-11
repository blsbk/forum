package models

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func (m *PostModel) Insert(title, content, author string, categories []string) (int, error) {
	stmt := `INSERT INTO posts (title, content, created, author, likes)
	VALUES(?, ?, datetime('now', 'utc'), ?, "0");`

	result, err := m.DB.Exec(stmt, title, content, author)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	err = m.CategoryInsert(id, categories)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *PostModel) CategoryInsert(postid int64, categories []string) error {
	stmt2 := `INSERT INTO categories (postid, category) VALUES (?, ?)`

	for _, category := range categories {
		_, err := m.DB.Exec(stmt2, postid, category)
		if err != nil {
			return err
		}
	}
	return nil
}

// This will insert a new post into the database.
func (m *PostModel) Get(id int) (*Post, error) {
	p := &Post{}

	stmt := `SELECT id, title, content, created, author, likes FROM posts WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Author, &p.Likes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return p, nil
}

// This will return the 10 most recently created posts.

func (m *PostModel) Latest() (map[int]*Post, error) {
	stmt := `SELECT id, title, content, created, author, likes FROM posts ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
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

func (m *PostModel) FilteredPosts(categories []string) (map[int]*Post, error) {
	stmt := `SELECT id, title, content, created, author, likes FROM posts JOIN categories ON posts.id = categories.postid WHERE categories.category = ?;`
	posts := map[int]*Post{}

	for _, category := range categories {
		rows, err := m.DB.Query(stmt, category)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			p := &Post{}

			err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Author, &p.Likes)
			if err != nil {
				return nil, err
			}

			posts[p.ID] = p
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func (m *PostModel) GetUserPosts(author string) (map[int]*Post, error) {
	stmt := `SELECT id, title, content, created, author, likes FROM posts WHERE author = ?`

	rows, err := m.DB.Query(stmt, author)
	if err != nil {
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

func (m *PostModel) LikeInsert(likeData UserLikeData, likedBy string) error {
	stmt := `UPDATE posts SET likes = ? WHERE ? = id`

	_, err := m.DB.Exec(stmt, likeData.Likes, likeData.ID)
	if err != nil {
		return err
	}

	var stmt2 string

	if likeData.IsLiked {
		stmt2 = `INSERT INTO liked (postid, by)
		VALUES(?, ?);`
		_, err = m.DB.Exec(stmt2, likeData.ID, likedBy)
		if err != nil {
			return err
		}
	} else {
		stmt2 = `DELETE FROM liked WHERE by = ?`
		_, err = m.DB.Exec(stmt2, likedBy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PostModel) IsLikedByUser(user string, postid int) bool {
	stmt := `SELECT EXISTS (SELECT * FROM liked WHERE by = ? AND postid = ?)`

	var exists bool

	err := m.DB.QueryRow(stmt, user, postid).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (m *PostModel) CommentInsert(comment, commentBy string, postId int) error {
	stmt := `INSERT INTO comments (postid, comment, by, likes) VALUES(?, ?, ?, '0');`

	_, err := m.DB.Exec(stmt, postId, comment, commentBy)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostModel) GetComments(postId int) ([]*PostComments, error) {
	stmt := `SELECT id, comment, by, likes FROM comments WHERE postid = ?;`
	rows, err := m.DB.Query(stmt, postId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*PostComments{}

	for rows.Next() {

		c := &PostComments{}

		err = rows.Scan(&c.Id, &c.Comment, &c.Author, &c.Likes)

		if err != nil {
			return nil, err
		}

		comments = append(comments, c)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *PostModel) CommentLikeInsert(likeData CommentLikeData, likedBy string) error {
	stmt := `UPDATE comments SET likes = ? WHERE ? = id`

	_, err := m.DB.Exec(stmt, likeData.Likes, likeData.ID)
	if err != nil {
		return err
	}

	var stmt2 string

	if likeData.IsLiked {
		stmt2 = `INSERT INTO comment_likes (commentid, postid, by)
		VALUES(?, ?, ?);`
		_, err = m.DB.Exec(stmt2, likeData.ID, likeData.PostId, likedBy)
		if err != nil {
			return err
		}
	} else {
		stmt2 = `DELETE FROM comment_likes WHERE by = ?`
		_, err = m.DB.Exec(stmt2, likedBy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PostModel) IsCommentLikedByUser(user string, commentid int) bool {
	stmt := `SELECT EXISTS (SELECT * FROM comment_likes WHERE by = ? AND commentid = ?)`

	var exists bool

	err := m.DB.QueryRow(stmt, user, commentid).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (m *PostModel) GetPostId(r *http.Request) (int, error) {
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
