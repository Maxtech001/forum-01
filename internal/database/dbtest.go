package database

import "fmt"

func TestDbStuff() {
	fmt.Println(DbGetPosts())
	// add post
	//DbInsertPost(user_id, title, content, tags)

	//add comment
	//DbInsertComment(2, "user1", "OMG you're totally right!")
	// register new user
	/*
		var newuser User
		newuser.Id = "user1"
		newuser.Name = "User 1"
		newuser.Email = "user.1@example.com"
		newuser.Password = "12345"

		err := DbInsertUser(newuser)

		if err != nil {
			fmt.Println("500", err)
		}
			// get user data by username or e-mail address, returns user struct
			//fmt.Println(dbGetUserByIdOrEmail("n00bh4ck3r"))
			fmt.Println(dbGetUserByIdOrEmail("margus.axx@gmail.com"))

			// authenticate user by email/username and password, returns boolean
			//fmt.Println(dbAuthenticateUser("n00bh4ck3r", "qwerty"))
			fmt.Println(dbAuthenticateUser("margusaid", "12345"))
	*/

}
