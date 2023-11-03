package repository

import (
	"database/sql"
	"errors"
	"strings"

	"forum.bbilisbe/internal/models"
)

type sqlPostsRepository struct {
	Conn *sql.DB
}

func NewSqlPostsRepository(conn *sql.DB) models.PostRepository {
	return &sqlPostsRepository{conn}
}

func (m *sqlPostsRepository) Insert(title, content, author string, categories []string) (int, error) {
	stmt := `INSERT INTO posts (title, content, created, author, likes, dislikes, tags)
	VALUES(?, ?, datetime('now', 'utc'), ?, "0", "0", ?);`

	result, err := m.Conn.Exec(stmt, title, content, author, strings.Join(categories, " "))
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

func (m *sqlPostsRepository) CategoryInsert(postid int64, categories []string) error {
	stmt2 := `INSERT INTO categories (postid, category) VALUES (?, ?)`

	for _, category := range categories {
		_, err := m.Conn.Exec(stmt2, postid, category)
		if err != nil {
			return err
		}
	}
	return nil
}

// This will insert a new post into the database.
func (m *sqlPostsRepository) Get(id int) (*models.Post, error) {
	p := &models.Post{}

	stmt := `SELECT id, title, content, created, author, likes, dislikes, tags FROM posts WHERE id = ?`

	err := m.Conn.QueryRow(stmt, id).Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Author, &p.Likes, &p.Dislikes, &p.Tags)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return p, nil
}

// This will return the 10 most recently created posts.

func (m *sqlPostsRepository) Latest() (map[int]*models.Post, error) {
	stmt := `SELECT id, title, author, likes FROM posts ORDER BY id DESC LIMIT 10`

	rows, err := m.Conn.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := map[int]*models.Post{}

	for rows.Next() {
		p := &models.Post{}

		err = rows.Scan(&p.ID, &p.Title, &p.Author, &p.Likes)
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

func (m *sqlPostsRepository) FilteredPosts(categories []string) (map[int]*models.Post, error) {
	posts := map[int]*models.Post{}
	stmt := `SELECT id, title, author, likes, tags FROM posts JOIN categories ON posts.id = categories.postid WHERE categories.category = ?;`

	for _, category := range categories {
		rows, err := m.Conn.Query(stmt, category)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			p := &models.Post{}

			err := rows.Scan(&p.ID, &p.Title, &p.Author, &p.Likes, &p.Tags)
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

func (m *sqlPostsRepository) LikeInsert(likeData models.UserLikeData, likedBy string) error {
	stmt := `UPDATE posts SET likes = ? WHERE ? = id`

	_, err := m.Conn.Exec(stmt, likeData.Likes, likeData.ID)
	if err != nil {
		return err
	}

	var stmt2 string

	if likeData.IsLiked {
		stmt2 = `INSERT INTO likes (postid, likedby)
		VALUES(?, ?);`
		_, err = m.Conn.Exec(stmt2, likeData.ID, likedBy)
		if err != nil {
			return err
		}
		if m.IsDislikedByUser(likedBy, likeData.ID) {
			m.RemoveDislike(likeData.ID, likedBy)
		}
	} else {
		stmt2 = `DELETE FROM likes WHERE likedby = ? AND postid = ?`
		_, err = m.Conn.Exec(stmt2, likedBy, likeData.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *sqlPostsRepository) RemoveLike(postid int, likedBy string) error {
	stmt := `DELETE FROM likes WHERE likedby = ? AND postid = ?`
	_, err := m.Conn.Exec(stmt, likedBy, postid)
	if err != nil {
		return err
	}
	var likes int
	stmt2 := `SELECT likes FROM posts WHERE id = ?`
	row := m.Conn.QueryRow(stmt2, postid)
	row.Scan(&likes)
	likes--

	stmt3 := `UPDATE posts SET likes = ? WHERE ? = id`

	_, err = m.Conn.Exec(stmt3, likes, postid)
	if err != nil {
		return err
	}
	return nil
}

func (m *sqlPostsRepository) DislikeInsert(dislikeData models.UserDislikeData, dislikedBy string) error {
	stmt := `UPDATE posts SET dislikes = ? WHERE ? = id`

	_, err := m.Conn.Exec(stmt, dislikeData.Dislikes, dislikeData.ID)
	if err != nil {
		return err
	}

	var stmt2 string

	if dislikeData.IsDisliked {
		stmt2 = `INSERT INTO dislikes (postid, dislikedby)
		VALUES(?, ?);`
		_, err = m.Conn.Exec(stmt2, dislikeData.ID, dislikedBy)
		if err != nil {
			return err
		}

		if m.IsLikedByUser(dislikedBy, dislikeData.ID) {
			m.RemoveLike(dislikeData.ID, dislikedBy)
		}
	} else {
		stmt2 = `DELETE FROM dislikes WHERE dislikedby = ? AND postid = ?`
		_, err = m.Conn.Exec(stmt2, dislikedBy, dislikeData.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *sqlPostsRepository) RemoveDislike(postid int, likedBy string) error {
	stmt := `DELETE FROM dislikes WHERE dislikedby = ? AND postid = ?`
	_, err := m.Conn.Exec(stmt, likedBy, postid)
	if err != nil {
		return err
	}
	var dislikes int
	stmt2 := `SELECT dislikes FROM posts WHERE id = ?`
	row := m.Conn.QueryRow(stmt2, postid)
	row.Scan(&dislikes)
	dislikes--

	stmt3 := `UPDATE posts SET dislikes = ? WHERE ? = id`

	_, err = m.Conn.Exec(stmt3, dislikes, postid)
	if err != nil {
		return err
	}
	return nil
}

func (m *sqlPostsRepository) IsLikedByUser(user string, postid int) bool {
	stmt := `SELECT EXISTS (SELECT * FROM likes WHERE likedby = ? AND postid = ?)`

	var exists bool

	err := m.Conn.QueryRow(stmt, user, postid).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (m *sqlPostsRepository) IsDislikedByUser(user string, postid int) bool {
	stmt := `SELECT EXISTS (SELECT * FROM dislikes WHERE dislikedby = ? AND postid = ?)`

	var exists bool

	err := m.Conn.QueryRow(stmt, user, postid).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (m *sqlPostsRepository) CommentInsert(comment, commentBy string, postId int) error {
	stmt := `INSERT INTO comments (postid, comment, commentby, likes, dislikes) VALUES(?, ?, ?, '0', '0');`

	_, err := m.Conn.Exec(stmt, postId, comment, commentBy)
	if err != nil {
		return err
	}

	return nil
}

func (m *sqlPostsRepository) GetComments(postId int, user string) ([]*models.PostComments, error) {
	stmt := `SELECT id, comment, commentby, likes, dislikes FROM comments WHERE postid = ?;`
	rows, err := m.Conn.Query(stmt, postId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*models.PostComments{}

	for rows.Next() {

		c := &models.PostComments{}

		err = rows.Scan(&c.Id, &c.Comment, &c.Author, &c.Likes, &c.Dislikes)

		if err != nil {
			return nil, err
		}
		c.IsLiked = m.IsCommentLikedByUser(user, c.Id)

		comments = append(comments, c)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *sqlPostsRepository) CommentLikeInsert(likeData models.CommentLikeData, likedBy string) error {
	stmt := `UPDATE comments SET likes = ? WHERE id = ?`

	_, err := m.Conn.Exec(stmt, likeData.Likes, likeData.ID)
	if err != nil {
		return err
	}

	var stmt2 string

	if likeData.IsLiked {
		stmt2 = `INSERT INTO comment_likes (commentid, postid, likedby)
		VALUES(?, ?, ?);`
		_, err = m.Conn.Exec(stmt2, likeData.ID, likeData.PostId, likedBy)
		if err != nil {
			return err
		}
		if m.IsCommentDislikedByUser(likedBy, likeData.ID) {
			m.RemoveCommentDislike(likeData.ID, likedBy)
		}
	} else {
		stmt2 = `DELETE FROM comment_likes WHERE commentid = ? AND likedby = ?`
		_, err = m.Conn.Exec(stmt2, likeData.ID, likedBy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *sqlPostsRepository) RemoveCommentLike(commentid int, likedBy string) error {
	stmt := `DELETE FROM comment_likes WHERE likedby = ? AND commentid = ?`
	_, err := m.Conn.Exec(stmt, likedBy, commentid)
	if err != nil {
		return err
	}
	var likes int
	stmt2 := `SELECT likes FROM comments WHERE id = ?`
	row := m.Conn.QueryRow(stmt2, commentid)
	row.Scan(&likes)
	likes--

	stmt3 := `UPDATE comments SET likes = ? WHERE ? = id`

	_, err = m.Conn.Exec(stmt3, likes, commentid)
	if err != nil {
		return err
	}
	return nil
}

func (m *sqlPostsRepository) IsCommentLikedByUser(user string, commentid int) bool {
	stmt := `SELECT EXISTS (SELECT * FROM comment_likes WHERE likedby = ? AND commentid = ?)`

	var exists bool

	err := m.Conn.QueryRow(stmt, user, commentid).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (m *sqlPostsRepository) CommentDislikeInsert(dislikeData models.CommentDislikeData, dislikedBy string) error {
	stmt := `UPDATE comments SET dislikes = ? WHERE id = ?`

	_, err := m.Conn.Exec(stmt, dislikeData.Dislikes, dislikeData.ID)
	if err != nil {
		return err
	}

	var stmt2 string

	if dislikeData.IsDisliked {
		stmt2 = `INSERT INTO comment_dislikes (commentid, postid, dislikedby)
		VALUES(?, ?, ?);`
		_, err = m.Conn.Exec(stmt2, dislikeData.ID, dislikeData.PostId, dislikedBy)
		if err != nil {
			return err
		}
		if m.IsCommentLikedByUser(dislikedBy, dislikeData.ID) {
			m.RemoveCommentLike(dislikeData.ID, dislikedBy)
		}
	} else {
		stmt2 = `DELETE FROM comment_dislikes WHERE commentid = ? AND dislikedby = ?`
		_, err = m.Conn.Exec(stmt2, dislikeData.ID, dislikedBy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *sqlPostsRepository) RemoveCommentDislike(commentid int, likedBy string) error {
	stmt := `DELETE FROM comment_dislikes WHERE dislikedby = ? AND commentid = ?`
	_, err := m.Conn.Exec(stmt, likedBy, commentid)
	if err != nil {
		return err
	}
	var dislikes int
	stmt2 := `SELECT dislikes FROM comments WHERE id = ?`
	row := m.Conn.QueryRow(stmt2, commentid)
	row.Scan(&dislikes)
	dislikes--

	stmt3 := `UPDATE comments SET dislikes = ? WHERE ? = id`

	_, err = m.Conn.Exec(stmt3, dislikes, commentid)
	if err != nil {
		return err
	}
	return nil
}

func (m *sqlPostsRepository) IsCommentDislikedByUser(user string, commentid int) bool {
	stmt := `SELECT EXISTS (SELECT * FROM comment_dislikes WHERE dislikedby = ? AND commentid = ?)`

	var exists bool

	err := m.Conn.QueryRow(stmt, user, commentid).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (m *sqlPostsRepository) GetPostId() {
	return
}
