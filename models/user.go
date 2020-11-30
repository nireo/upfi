package models

import (
	"os"

	"github.com/nireo/upfi/lib"
	"gorm.io/gorm"
)

// User is a database struct, which also holds all the properties of gorm.Model
type User struct {
	gorm.Model
	Username string
	Password string
	UUID     string
	Files    []File
}

// Serialize serializes a given user's data into json format
func (user *User) Serialize() lib.JSON {
	return lib.JSON{
		"username": user.Username,
		"uuid":     user.UUID,
	}
}

// Delete deletes the given user's file and removes the user's database entry from the database.
func (user *User) Delete() error {
	db := lib.GetDatabase()

	// Remove the user's folder
	err := os.RemoveAll("./files/" + user.UUID)
	if err != nil {
		return err
	}

	// Remove from database
	db.Delete(&user)
	return nil
}

// FindOneUser takes a interface{} as an argument and returns a pointer to a user struct,
// if it is found.
func FindOneUser(condition interface{}) (*User, error) {
	db := lib.GetDatabase()
	var user User
	if err := db.Where(condition).First(&user).Error; err != nil {
		// Not found so returns a null pointer.
		return nil, err
	}

	return &user, nil
}
