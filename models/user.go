package models

import (
	"os"

	"github.com/nireo/upfi/lib"
	"gorm.io/gorm"
)

// User is a database struct, which also holds all the properties of gorm.Model
type User struct {
	gorm.Model
	Username             string `json:"username"` // Username used to login.
	Password             string // Password to see the files.
	UUID                 string `json:"uuid"` // Unique ID to identify a user.
	FileEncryptionMaster string // A password which holds the passphrase with which files are encrypted.
	Files                []File // A relation to files, which hold a UserID which refers to this model.
}

// Serialize serializes a given user's data into json format
func (user *User) Serialize() lib.JSON {
	return lib.JSON{
		"username": user.Username,
		"uuid":     user.UUID,
	}
}

func (user *User) FindSharedToFiles() ([]File, error) {
	db := lib.GetDatabase()

	var sharedFileContracts []FileShare

	db.Where(&FileShare{SharedToID: user.ID}).Find(&sharedFileContracts)

	var files []File
	for _, sf := range sharedFileContracts {
		var tempFile File
		if err := db.Where("id = ?", sf.SharedFileID).First(&tempFile).Error; err != nil {
			return nil, err
		}

		files = append(files, tempFile)
	}

	return files, nil
}

func (user *User) FindSharedByFiles() ([]File, error) {
	db := lib.GetDatabase()

	var sharedFileContracts []FileShare
	db.Where(&FileShare{SharedByID: user.ID}).Find(&sharedFileContracts)

	var files []File
	for _, sf := range sharedFileContracts {
		var tempFile File
		if err := db.Where("id = ?", sf.SharedFileID).First(&tempFile).Error; err != nil {
			return nil, err
		}

		files = append(files, tempFile)
	}

	return files, nil
}

// Delete deletes the given user's file and removes the user's database entry from the database.
func (user *User) Delete() error {
	db := lib.GetDatabase()

	// Remove the user's folder
	if err := os.RemoveAll(lib.AddRootToPath("files/") + user.UUID); err != nil {
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
