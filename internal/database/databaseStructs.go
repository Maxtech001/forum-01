package database

type User struct {
	Id       string
	Name     string
	Email    string
	Password string
}

type Post struct {
	Id       int
	User_id  string
	Title    string
	Time     string
	Content  string
	Comments []Comment
	Likes    int
	Dislikes int
	Tags     []Tag
}

type Comment struct {
	Id       int
	User_id  string
	Time     string
	Content  string
	Likes    int
	Dislikes int
}

type Tag struct {
	Id   int
	Name string
}

type Mainpage struct {
	Username string
	Posts    []Post
	Tags     []Tag
}
