package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type User struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type UserDB struct {
	Path  string
	Users []*User
}

func New(p string) (db *UserDB, err error) {

	db = &UserDB{Path: p}

	err = db.LoadUsers()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *UserDB) LoadUsers() error {

	// Check if the file does not exist yet and if it doesn't, create one.
	if _, err := os.Stat(db.Path); os.IsNotExist(err) {
		db.createUsersDBFIle()
	}

	// Read out the contents of the file.
	buf, err := ioutil.ReadFile(db.Path)
	if err != nil {
		return err
	}

	// Serialize the file data into the list of users.
	err = json.Unmarshal(buf, &db.Users)
	if err != nil {
		return err
	}

	return nil
}

func (db *UserDB) CreateUser(name string, code string) (user *User, err error) {

	// Make sure our list of users is current.
	if err = db.LoadUsers(); err != nil {
		return nil, err
	}

	// TODO: Check uniqueness here...
	exists := db.userExists(code)
	if exists {
		return nil, fmt.Errorf("User %s already exists", name)
	}

	user = &User{name, code}
	db.Users = append(db.Users, user)

	buf, err := json.Marshal(db.Users)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(db.Path, buf, 0600)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *UserDB) FindUser(code string) (user *User, err error) {

	// Make sure our list of users is current.
	if err = db.LoadUsers(); err != nil {
		return nil, err
	}

	// Check to see if the user already exists in our system.
	for _, u := range db.Users {
		if u.Code == code {
			return u, nil
		}
	}

	return nil, fmt.Errorf("No user could be found with the code '%s'", code)
}

func (db *UserDB) createUsersDBFIle() error {

	err := ioutil.WriteFile(db.Path, []byte("[]"), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (db *UserDB) userExists(code string) (exists bool) {

	// Make sure our list of users is current.
	if err := db.LoadUsers(); err != nil {
		return false
	}

	// Check to see if the user already exists in our system.
	for _, u := range db.Users {
		if u.Code == code {
			return true
		}
	}

	return false
}
