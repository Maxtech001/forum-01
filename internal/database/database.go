package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbfile = "./db/forum.db"

var Db *sql.DB

func InitDb() {
	// Setting up database
	Db = DbOpen()
}

// open database
func DbOpen() *sql.DB {
	Db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return Db
}

// get posts
func DbGetPosts() []Post {
	var result []Post
	sql := "select id, user_id, time, title, content, " +
		"(select count(*) from feedback f where f.post_id=p.id and f.type = '+') likes, " +
		"(select count(*) from feedback f where f.post_id=p.id and f.type = '-') dislikes, " +
		"(select count(*) from comment c where c.post_id=p.id) comments " +
		"from post p " +
		"order by time desc"
	rows, err := Db.Query(sql)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.User_id, &post.Time, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.Comments)
		if err != nil {
			fmt.Println(err)
			return result
		}
		post.Tags = DbGetPostTags(post.Id)
		result = append(result, post)
	}
	return result

}

// get post comments
func DbGetPostComments(post_id int) []Comment {
	var result []Comment
	sql := "select id, user_id, time, content, " +
		"(select count(*) from feedback f where f.comment_id=c.id and f.type = '+') likes," +
		"(select count(*) from feedback f where f.comment_id=c.id and f.type = '-') dislikes " +
		"from comment c where c.post_id=?"
	rows, err := Db.Query(sql, post_id)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.Id, &comment.User_id, &comment.Time, &comment.Content, &comment.Likes, &comment.Dislikes)
		if err != nil {
			fmt.Println(err)
			return result
		}
		result = append(result, comment)
	}
	return result
}

// get all tags
func DbGetTags() []Tag {
	var result []Tag
	rows, err := Db.Query("SELECT id, name FROM tag")
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		err = rows.Scan(&tag.Id, &tag.Name)
		if err != nil {
			fmt.Println(err)
			return result
		}
		result = append(result, tag)
	}
	return result
}

// get post tags
func DbGetPostTags(post_id int) []Tag {
	var result []Tag
	rows, err := Db.Query("SELECT pt.tag_id, t.name FROM post_tag pt LEFT JOIN tag t ON pt.tag_id = t.id WHERE pt.post_id=?", post_id)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		err = rows.Scan(&tag.Id, &tag.Name)
		if err != nil {
			fmt.Println(err)
			return result
		}
		result = append(result, tag)
	}
	return result
}

// insert feedback
func DbInsertFeedback(post_id, comment_id int, user_id, ftype string) error {
	fbq, err := Db.Prepare("INSERT INTO feedback(post_id, comment_id, user_id, type) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer fbq.Close()
	_, err = fbq.Exec(post_id, comment_id, user_id, ftype)
	if err != nil {
		return err

	}

	return nil
}

// insert comment
func DbInsertComment(post_id int, user_id, content string) error {
	t := time.Now()
	dbtime := t.Format("2006-01-02 15:04:05")
	commentq, err := Db.Prepare("INSERT INTO comment(post_id, user_id, content, time) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer commentq.Close()

	_, err = commentq.Exec(post_id, user_id, content, dbtime)
	if err != nil {
		return err

	}

	return nil
}

// insert new post and its tags
func DbInsertPost(user_id, title, content string, tags []int) error {
	t := time.Now()
	dbtime := t.Format("2006-01-02 15:04:05")
	postq, err := Db.Prepare("INSERT INTO post(user_id, title, content, time) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer postq.Close()

	_, err = postq.Exec(user_id, title, content, dbtime)
	if err != nil {
		return err

	}
	post_id := 0
	err = Db.QueryRow("SELECT id FROM post WHERE time=? and title=? and content=? and user_id=?", dbtime, title, content, user_id).Scan(&post_id)
	if err != nil {
		return err
	}

	if post_id == 0 {
		return errors.New("could not find the post")
	}

	tagq, err := Db.Prepare("INSERT INTO post_tag(post_id, tag_id) values(?, ?)")
	if err != nil {
		return err
	}
	defer tagq.Close()

	for _, tag := range tags {

		_, err = tagq.Exec(post_id, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

// insert new user
func DbInsertUser(user User) error {
	stmt, err := Db.Prepare("INSERT INTO user(id, email, password) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	pwd := HashPassword(user.Password)

	_, err = stmt.Exec(user.Id, user.Email, pwd)
	if err != nil {
		return err

	}
	return nil
}

// check if user exists
func DbGetUserByIdOrEmail(input string) []User {
	var result []User
	rows, err := Db.Query("SELECT id, email, password FROM user WHERE id=? OR email=?", input, input)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Email, &user.Password)
		if err != nil {
			fmt.Println(err)
			return result
		}
		result = append(result, user)
	}
	return result
}

// check if user with email exist
func DbEmailExist(input string) bool {
	var usercount int

	err := Db.QueryRow("SELECT count(*) FROM user WHERE email=?", input).Scan(&usercount)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if usercount == 0 {
		return false
	} else {
		return true
	}
}

// check if user with isername exist
func DbUserIdExist(input string) bool {
	var usercount int

	err := Db.QueryRow("SELECT count(*) FROM user WHERE id=?", input).Scan(&usercount)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if usercount == 0 {
		return false
	} else {
		return true
	}
}

// authenticate by username and password
func DbAuthenticateUser(email, pwd string) bool {
	result := false
	var pw string

	err := Db.QueryRow("SELECT password FROM user WHERE email=?", email).Scan(&pw) //todo, use count(*) instead
	if err != nil {
		return result
	}
	return CheckPasswordHash(pwd, pw)
}
