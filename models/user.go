package models

import (
	"os"

	"github.com/nireo/upfi/lib"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
	UUID     string
	Files    []File
}

func (user *User) Serialize() lib.JSON {
	return lib.JSON{
		"username": user.Username,
		"uuid":     user.UUID,
	}
}

func (user *User) Delete() error {
	db := lib.GetDatabase()

	err := os.RemoveAll("./files/" + user.UUID)
	if err != nil {
		return err
	}

	db.Delete(&user)
	return nil
}

func FindOneUser(condition interface{}) (*User, error) {
	db := lib.GetDatabase()
	var user User
	if err := db.Where(condition).First(&user).Error; err != nil {
		return &user, err
	}

	return &user, nil
}
