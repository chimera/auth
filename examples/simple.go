package main

import (
	"fmt"
	"os"

	"github.com/chimera/auth"
)

func main() {

	userFile := "users.json"

	// Open up the user file, load any existing users and return them.
	db, err := auth.New(userFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a few users.
	fluffy, err := db.CreateUser("Fluffy", "cats")
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("Created user:", fluffy)

	spike, err := db.CreateUser("Spike", "roof")
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("Created user:", spike)

	// Search for a user in the list of users provided the given code.
	user, err := db.FindUser("cats")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Found user:", user)

	// Display all the users in the database.
	fmt.Println("All users:")
	for _, u := range db.Users {
		fmt.Printf("\tUser: %v\n", u)
	}

	// Remove the user file to start with a clean slate.
	os.Remove(userFile)

}
