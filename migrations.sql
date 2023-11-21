PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS posts(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		title VARCHAR(100) NOT NULL,
		content TEXT NOT NULL,
		created DATETIME NOT NULL,
		author NUMBER NOT NULL,
		likes NUMBER,
		dislikes NUMBER,
		tags TEXT NOT NULL,
		image TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_posts_created ON posts(created);


	CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	username VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	hashed_password CHAR(60) NOT NULL,
	token TEXT,
	expiry DATETIME,
	created DATETIME NOT NULL,
	CONSTRAINT unique_email UNIQUE (email),
	CONSTRAINT unique_name UNIQUE (username)
	);

	UPDATE users SET token = NULL, expiry = NULL;

	CREATE TABLE IF NOT EXISTS likes (
		postid INTEGER,
		likedby NUMBER
	);


	CREATE TABLE IF NOT EXISTS dislikes (
		postid INTEGER,
		dislikedby NUMBER
	);

	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		postid INTEGER,
		comment TEXT,
		likes INTEGER,
		dislikes INTEGER,
		commentby NUMBER
	);

	CREATE TABLE IF NOT EXISTS comment_likes (
		commentid INTEGER,
		postid INTEGER,
		likedby NUMBER
	);

	CREATE TABLE IF NOT EXISTS comment_dislikes (
		commentid INTEGER,
		postid INTEGER,
		dislikedby NUMBER
	);


	CREATE TABLE IF NOT EXISTS categories (
		postid INTEGER,
		category TEXT
	);
