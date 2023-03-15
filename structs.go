package main

type user struct {
	Id       string
	Name     string
	Email    string
	Password string
}

type post struct {
	Id       int
	User_id  int
	Title    string
	Time     string
	Content  string
	Comments []comment
	Likes    int
	Dislikes int
	Tags     []tag
}

type comment struct {
	Id       int
	User_id  string
	Time     string
	Content  string
	Likes    int
	Dislikes int
}

type tag struct {
	Id   int
	Name string
}

type mainpage struct {
	Username string
	Posts    []post
}
