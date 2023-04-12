package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbfile = "./db/forum.db"

var (
	db   *sql.DB
	once sync.Once
)

// Init db connection and only checking once
func InitDB() (*sql.DB, error) {
	var err error
	once.Do(func() {
		db, err = sql.Open("sqlite3", dbfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = db.Ping()
		if err != nil {
			fmt.Println(err)
			return
		}
	})
	return db, err
}

func DbGetUserByCookie(cookie string) string {

	if cookie == "" {
		return ""
	}
	var user string
	err := db.QueryRow("SELECT user_id FROM session WHERE id=? and datetime(expires) > datetime('now')", cookie).Scan(&user)
	if err != nil {
		return ""
	}

	if user != ""{
		dbq, _ := db.Prepare("DELETE FROM session WHERE user_id=? AND id<>?")
		defer dbq.Close()
		dbq.Exec(user,cookie)	
	}

	return user
}

func DbDeleteCookie(cookie string) {
	dbq, _ := db.Prepare("DELETE FROM session WHERE id = ?")

	defer dbq.Close()
	dbq.Exec(cookie)
}

func DbDeleteExpiredCookies() {
	stmt, _ := db.Prepare("DELETE FROM session WHERE datetime(expires) < datetime('now') or expires is NULL")
	defer stmt.Close()
	stmt.Exec()
}

func DbAddCookie(cookie, user_id string, exp time.Time) {
	if user_id == "" {
		fmt.Println("user_id missing, can't set cookie")
		return
	}
	dbtime := exp.Format("2006-01-02 15:04:05")

	dbq, _ := db.Prepare("INSERT INTO session(id, user_id, expires) values (?, ?, ?)")

	defer dbq.Close()
	dbq.Exec(cookie, user_id, dbtime)
}

// get single post
func DbGetSinglePost(post_id int, user_id string) Post {
	var result Post
	sql := "select id, user_id, time, title, content, " +
		"(select count(*) from feedback f where f.post_id=p.id and f.type = '+') likes, " +
		"(select count(*) from feedback f where f.post_id=p.id and f.type = '-') dislikes, " +
		"(select count(*) from feedback f where f.post_id=p.id and f.type = '+' and f.user_id='" + user_id + "') hasliked, " +
		"(select count(*) from feedback f where f.post_id=p.id and f.type = '-' and f.user_id='" + user_id + "') hasdisliked, " +
		"(select count(*) from comment c where c.post_id=p.id) comments " +
		"from post p " +
		"where id =?"
	rows, err := db.Query(sql, post_id)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.User_id, &post.Time, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.HasLiked, &post.HasDisliked, &post.Comments)
		if err != nil {
			fmt.Println(err)
			return result
		}
		post.Commentstruct = DbGetPostComments(post.Id, user_id)
		post.Tags = DbGetPostTags(post.Id)
		result = post
	}

	//	fmt.Println(result)
	return result
}

// get posts
func DbGetPosts(user_id string, params map[string][]string) []Post {
	var result []Post

	tagfilter := ""
	if len(params["tag"]) > 0 {
		for _, tag := range params["tag"] {
			if tagfilter == "" {
				tagfilter = "where ("
			} else {
				tagfilter = tagfilter + "or "
			}
			tagfilter = tagfilter + "(select count(*) from post_tag pt where pt.post_id=p.id and pt.tag_id=" + tag + ")>0 "
		}
		tagfilter = tagfilter + ") "
	}

	ownpostfilter := ""
	if len(params["ownposts"]) > 0 {
		if tagfilter == "" {
			ownpostfilter = "where "
		} else {
			ownpostfilter = "and "
		}
		ownpostfilter = ownpostfilter + "user_id ='" + user_id + "'"
	}

	likefilter := ""
	if len(params["liked"]) > 0 {
		if tagfilter == "" && ownpostfilter == "" {
			likefilter = "where "
		} else {
			likefilter = "and "
		}
		likefilter = likefilter + "(select count(*) from feedback f where f.post_id = p.id and user_id ='" + user_id + "')>0 "
	}

	sql := "select id, user_id, time, title, content, " +
		"(select count(*) from feedback f where f.post_id=p.id and f.type = '+') likes, " +
		"(select count(*) from feedback f where f.post_id=p.id and f.type = '-') dislikes, " +
		"(select count(*) from feedback f where f.post_id=p.id and f.type = '+' and f.user_id='" + user_id + "') hasliked, " +
		"(select count(*) from feedback f where f.post_id=p.id and f.type = '-' and f.user_id='" + user_id + "') hasdisliked, " +
		"(select count(*) from comment c where c.post_id=p.id) comments " +
		"from post p " +
		tagfilter +
		ownpostfilter +
		likefilter +
		"order by time desc"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.User_id, &post.Time, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.HasLiked, &post.HasDisliked, &post.Comments)
		if err != nil {
			fmt.Println(err)
			return result
		}
		post.Commentstruct = DbGetPostComments(post.Id, user_id)
		post.Tags = DbGetPostTags(post.Id)
		result = append(result, post)
	}
	return result

}

// get post comments
func DbGetPostComments(post_id int, user_id string) []Comment {
	var result []Comment
	sql := "select id, post_id, user_id, time, content, " +
		"(select count(*) from feedback f where f.comment_id=c.id and f.type = '+') likes," +
		"(select count(*) from feedback f where f.comment_id=c.id and f.type = '-') dislikes, " +
		"(select count(*) from feedback f where f.comment_id=c.id and f.type = '+' and f.user_id='" + user_id + "') hasliked, " +
		"(select count(*) from feedback f where f.comment_id=c.id and f.type = '-' and f.user_id='" + user_id + "') hasdisliked " +
		"from comment c where c.post_id=?"
	rows, err := db.Query(sql, post_id)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.Id, &comment.Post_id, &comment.User_id, &comment.Time, &comment.Content, &comment.Likes, &comment.Dislikes, &comment.HasLiked, &comment.HasDisliked)
		if err != nil {
			fmt.Println(err)
			return result
		}
		comment.Post_id = post_id
		result = append(result, comment)
	}
	return result
}

// get all tags
func DbGetTags() []Tag {
	var result []Tag
	rows, err := db.Query("SELECT id, name FROM tag")
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
	rows, err := db.Query("SELECT pt.tag_id, t.name FROM post_tag pt LEFT JOIN tag t ON pt.tag_id = t.id WHERE pt.post_id=?", post_id)
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
	fbq, err := db.Prepare("INSERT INTO feedback(post_id, comment_id, user_id, type) values(?, ?, ?, ?)")
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
	commentq, err := db.Prepare("INSERT INTO comment(post_id, user_id, content, time) values(?, ?, ?, ?)")
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
func DbInsertPost(user_id, title, content string, tags []int) (error, int) {
	t := time.Now()
	dbtime := t.Format("2006-01-02 15:04:05")
	postq, err := db.Prepare("INSERT INTO post(user_id, title, content, time) values(?, ?, ?, ?)")
	if err != nil {
		return err, 0
	}
	defer postq.Close()

	_, err = postq.Exec(user_id, title, content, dbtime)
	if err != nil {
		return err, 0

	}
	post_id := 0
	err = db.QueryRow("SELECT id FROM post WHERE time=? and title=? and content=? and user_id=?", dbtime, title, content, user_id).Scan(&post_id)
	if err != nil {
		return err, 0
	}

	if post_id == 0 {
		return errors.New("could not find the post"), 0
	}

	tagq, err := db.Prepare("INSERT INTO post_tag(post_id, tag_id) values(?, ?)")
	if err != nil {
		return err, 0
	}
	defer tagq.Close()

	for _, tag := range tags {
		_, err = tagq.Exec(post_id, tag)
		if err != nil {
			return err, 0
		}
	}

	return nil, post_id
}

// insert new user
func DbInsertUser(user User) error {
	stmt, err := db.Prepare("INSERT INTO user(id, email, password) values(?, ?, ?)")
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
	rows, err := db.Query("SELECT id, email, password FROM user WHERE id=? OR email=?", input, input)
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

	err := db.QueryRow("SELECT count(*) FROM user WHERE email=?", input).Scan(&usercount)
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

	err := db.QueryRow("SELECT count(*) FROM user WHERE id=?", input).Scan(&usercount)
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

func DbAuthenticateUser(email, pwd string) (bool, string) {
	result := false
	var user, pw string

	err := db.QueryRow("SELECT id, password FROM user WHERE email=?", email).Scan(&user, &pw)
	if err != nil {
		return result, ""
	}

	if CheckPasswordHash(pwd, pw) {
		result = true
	} else {
		result = false
	}

	return result, user
}
