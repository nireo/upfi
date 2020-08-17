package models

import (
	"github.com/jinzhu/gorm"
	"github.com/nireo/upfi/lib"
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
		"uuid": user.UUID,
	}
}

func (user *User) Delete() error {
	db := lib.GetDatabase()

	err := os.RemoveAll("./files/"+user.UUID)
	if err != nil {
		return err
	}

	err := db.Delete(&user)
	if err != nil {
		return err
	}

	return nil
}

func FindOneUser(condition interface{}) (*User, error) {
	db := common.GetDatabase()
	var user models.User
	if err := db.Where(condition).First(&user).Error; err != nil {
		return &user, err
	}

	return &user, nil
}
