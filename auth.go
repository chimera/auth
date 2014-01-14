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

// func New(p string) error {
// 	path = p

// 	err := loadUsers()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (users *Users) Update() (err error) {
// 	err = loadUsers(users)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// Find searches for a user in the users file given a specified code.
// func FindUser(code string) (user *User, err error) {
// 	for _, u := range Users {
// 		if u.Code == code {
// 			return u, nil
// 		}
// 	}

// 	return nil, fmt.Errorf("No user could be found with the code '%s'", code)
// }

// func CreateUser(name string, code string) (user *User, err error) {

// 	// Open up the users file.
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		err := file.Close()
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()

// 	// Create the new users and append them to the list of users. We cache
// 	// the original set of users in case we fail at updating the users file.
// 	user = &User{name, code}
// 	// origUsers := users.All
// 	Users = append(Users, user)

// 	// Marshal all the users.
// 	buf, err := json.Marshal(Users)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Write the updated list of users to the users file.
// 	_, err = file.Write(buf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Attempt to update the users file with the new data.
// 	// err = users.updateUserFile()
// 	// if err != nil {
// 	// 	users.All = origUsers
// 	// 	return nil, err
// 	// }

// 	return user, nil
// }

// loadUsers opens a JSON users file and returns a Users struct populated with
// what it finds in the file, or an error.
// func loadUsers() error {

// 	file, err := openUserFile(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		err := file.Close()
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()

// 	// Read file contents into a list of bytes
// 	file_content := bytes.NewBuffer(nil)
// 	_, err = io.Copy(file_content, file)
// 	if err != nil {
// 		return err
// 	}

// 	// Unmarshall the file content into a slice of User structs.
// 	err = json.Unmarshal(file_content.Bytes(), &Users)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func openUserFile(p string) (file *os.File, err error) {

// 	// Check if the file does not exist yet and if it doesn't, create one.
// 	if _, err := os.Stat(p); os.IsNotExist(err) {

// 		// Create file.
// 		file, err := os.Create(p)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Make the file only readable/writable by the current user.
// 		err = file.Chmod(0600)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Add an empty JSON list to the file so it can be parsed by the JSON marshaller.
// 		_, err = file.Write([]byte("[]"))
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {

// 		// Open up the file for writing.
// 		file, err = os.Open(p)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return file, nil
// }

// func (users *Users) updateUserFile() (err error) {

// 	// Open up the users file.
// 	file, err := os.Open(users.Path) // users.openUserFile()
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		err = file.Close()
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()

// 	// Marshal all the users.
// 	buf, err := json.Marshal(users.All)
// 	if err != nil {
// 		return err
// 	}

// 	// Write the updated list of users to the users file.
// 	_, err = file.Write(buf)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// import (
//   "encoding/json"
//   "fmt"
//   "io/ioutil"
//   "log"
//   "net/http"
//   "os"

//   "github.com/codegangsta/martini-contrib/binding"
// )

// type UserDB struct {
//   Path string
// }

// type User struct {
//   Name string `json:"name" binding:"required"`
//   Code string `json:"code" binding:"required"`
// }

// // Create the Users database file.
// func (db *UserDB) Create() (err error) {
//   _, err = db.Write([]byte("[]"))
//   if err != nil {
//     return err
//   }
//   return nil
// }

// func (db *UserDB) Write(buf []byte) (n int, err error) {

//   // Create file.
//   file, err := os.Create(db.Path)
//   if err != nil {
//     return 0, err
//   }
//   defer file.Close()

//   // Make the file only readable/writable by the current user.
//   err = os.Chmod(db.Path, 0600)
//   if err != nil {
//     return 0, err
//   }

//   // Add an empty JSON list to the file so it can be parsed by the JSON marshaller.
//   n, err = file.Write(buf)
//   if err != nil {
//     return 0, err
//   }

//   return n, nil
// }

// func (db *UserDB) Read() (users []User, err error) {

//   // Check if the file does not exist yet and if it doesn't, create one.
//   if _, err := os.Stat(db.Path); os.IsNotExist(err) {

//     // Let user know we created a user file for them.
//     log.Print("No users JSON file, creating one now: ", db.Path)

//     err = db.Create()
//     if err != nil {
//       return nil, err
//     }
//   }

//   // Read file contents into a list of bytes
//   bytes, err := ioutil.ReadFile(db.Path)
//   if err != nil {
//     return nil, err
//   }

//   // Unmarshall the file content into a slice of User structs.
//   err = json.Unmarshal(bytes, &users)
//   if err != nil {
//     return nil, err
//   }

//   return users, nil
// }

// func (db *UserDB) CreateUser(u User) (user User, err error) {

//   // Get an updated list of users in the DB.
//   users, err := db.Read()
//   if err != nil {
//     return User{}, err
//   }

//   if userExists(u, users) {
//     return User{}, fmt.Errorf("%s already exists in the database", u.Name)
//   }
//   // Append the users to the list of users and write to the file.
//   users = append(users, u)
//   bytes, err := json.Marshal(users)
//   if err != nil {
//     return User{}, err
//   }

//   // Write the users back into the file
//   _, err = db.Write(bytes)
//   if err != nil {
//     return User{}, err
//   }

//   return u, nil
// }

// func (db *UserDB) AuthenticateCode(code string) (err error) {

//   // Make sure the list of users is up to date.
//   users, err := db.Read()
//   if err != nil {
//     return err
//   }

//   // Loop through all the users to see if any have a matching code.
//   for _, user := range users {
//     if code == user.Code {
//       logAccess(user)
//       return nil
//     }
//   }

//   // If the code doesn't match an entry in the database, return an error.
//   return fmt.Errorf("Your code '%s' is invalid, please try again!", code)
// }
