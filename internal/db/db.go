package dbase

import "database/sql"

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Exec(db *sql.DB) error {
	sts := `
	DROP TABLE IF EXISTS posts;

	CREATE TABLE posts(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		title VARCHAR(100) NOT NULL,
		content TEXT NOT NULL,
		created DATETIME NOT NULL,
		author TEXT,
		likes NUMBER
	);

	CREATE INDEX idx_posts_created ON posts(created);

	INSERT INTO posts (title, content, created, author, likes) VALUES (
	'New Dorama "Tomorrow"',
	'KBS is releasing new drama starring Idol-Actor from SF9 Rowoon',
	datetime('now', 'utc'),
	'Bagdat',
	'0'
	);
	INSERT INTO posts (title, content, created, author, likes) VALUES (
	'STRAY KIDS new release!',
	'STRAY KIDS from JYPE is releasing new song called "CASE 134". 
	Breaking the records on Bilboard-100.',
	datetime('now', 'utc'),
	'Bagdat',
	'0'
	);
	INSERT INTO posts (title, content, created, author, likes) VALUES (
	'BLACKPINK Lisa is dating someone?',
	'Member of BLACKPINK Lisa is romoured to be dating
	 a Hollywood movie star. YG Entertainment is staying silent 
	 on this situation, while Lisa herself is going on a vacation 
	 with his family.',
	datetime('now', 'utc'),
	'Yerassyl',
	'0'
	);
	
	DROP TABLE IF EXISTS users;

	CREATE TABLE users (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	username VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	hashed_password CHAR(60) NOT NULL,
	token TEXT,
	expiry DATETIME,
	created DATETIME NOT NULL,
	CONSTRAINT unique_email UNIQUE (email)
	);



	DROP TABLE IF EXISTS liked;

	CREATE TABLE liked (
		postid INTEGER PRIMARY KEY,
		by TEXT
	);

	DROP TABLE IF EXISTS comments;

	CREATE TABLE comments (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		postid INTEGER,
		comment TEXT,
		likes INTEGER,
		by TEXT
	);

	DROP TABLE IF EXISTS comment_likes;

	CREATE TABLE comment_likes (
		commentid INTEGER,
		postid INTEGER,
		by TEXT
	);

	DROP TABLE IF EXISTS categories;

	CREATE TABLE categories (
		postid INTEGER,
		category TEXT
	);

	INSERT INTO categories (postid, category) VALUES (
		'1',
		'dramas'
	);
	INSERT INTO categories (postid, category) VALUES (
		'1',
		'idols'
	);
	INSERT INTO categories (postid, category) VALUES (
		'2',
		'music'
	);
	INSERT INTO categories (postid, category) VALUES (
		'2',
		'idols'
	);
	INSERT INTO categories (postid, category) VALUES (
		'3',
		'idols'
	);
	`
	_, err := db.Exec(sts)
	if err != nil {
		return err
	}
	return nil
}
