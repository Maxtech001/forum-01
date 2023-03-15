package main

import "fmt"

func main() {
	db = dbOpen()

	if db == nil {
		fmt.Println("500")
	}
	defer db.Close()

	// register new user
	/*
		var newuser user
		newuser.Id = "margusaid"
		newuser.Name = "Margus Aid"
		newuser.Email = "margus.axx@gmail.com"
		newuser.Password = "12345"

		err := dbInsertUser(newuser)

		if err != nil {
			fmt.Println("500", err)
		}
	*/

	// get user data by username or e-mail address, returns user struct
	//fmt.Println(dbGetUserByIdOrEmail("n00bh4ck3r"))
	fmt.Println(dbGetUserByIdOrEmail("margus.axx@gmail.com"))

	// authenticate user by email/username and password, returns boolean
	//fmt.Println(dbAuthenticateUser("n00bh4ck3r", "qwerty"))
	fmt.Println(dbAuthenticateUser("margusaid", "12345"))

}
