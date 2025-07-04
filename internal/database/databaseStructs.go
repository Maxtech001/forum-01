package database

type User struct {
	Id       string
	Email    string
	Password string
}

type Post struct {
	Id            int
	User_id       string
	Title         string
	Time          string
	Content       string
	Comments      int
	Likes         int
	Dislikes      int
	HasLiked      int
	HasDisliked   int
	Tags          []Tag
	Commentstruct []Comment
}

type Comment struct {
	Id          int
	User_id     string
	Time        string
	Content     string
	Likes       int
	Dislikes    int
	HasLiked    int
	HasDisliked int
	Post_id     int
}

type Tag struct {
	Id   int
	Name string
}

type Mainpage struct {
	User_id string
	Posts   []Post
	Tags    []Tag
}

type Createpost struct {
	User_id string
	Tags    []Tag
}

type Postpage struct {
	User_id string
	Post    Post
}

type Aboutpage struct {
	User_id string
	Content string
}