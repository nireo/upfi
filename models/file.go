package models

import (
	"os"

	"github.com/nireo/upfi/lib"
	"gorm.io/gorm"
)

// File is a database struct, which also holds properties of gorm.Model
type File struct {
	gorm.Model
	Filename    string `json:"filename"`
	UUID        string `json:"uuid"`
	Description string `json:"description"`
	Size        int64  `json:"size"`
	SizeHuman   string
	UserID      uint
	Extension   string `json:"extension"`
	MIME        string `json:"mime"`
}

// FileShare represents a file share record
type FileShare struct {
	gorm.Model
	SharedByID   uint // the sharer's primary key
	SharedToID   uint // the primary key of the shared user
	SharedFileID uint
}

func (file *File) IsSharedTo(userID uint) bool {
	db := lib.GetDatabase()
	if err := db.Where(&FileShare{SharedToID: userID, SharedFileID: file.ID}).Error; err != nil {
		return false
	}

	return true
}

func (file *File) ShareToUser(shareTo *User) {
	db := lib.GetDatabase()

	fileShare := &FileShare{
		SharedByID:   file.UserID,
		SharedFileID: file.ID,
		SharedToID:   shareTo.ID,
	}

	db.Create(fileShare)
}

// Serialize serializes the user's data into json format
func (file *File) Serialize() lib.JSON {
	return lib.JSON{
		"filename":    file.Filename,
		"created_at":  file.CreatedAt,
		"description": file.Description,
		"uuid":        file.UUID,
	}
}

// Delete removes a given file and it's database entry
func (file *File) Delete(userID string) error {
	db := lib.GetDatabase()

	if err := os.Remove(lib.AddRootToPath("files/") + userID + "/" + file.Filename); err != nil {
		return err
	}

	db.Delete(&file)
	return nil
}

// FindOneFile takes a query interface{} as a parameter and returns a pointer to a file,
// if it is found.
func FindOneFile(condition interface{}) (*File, error) {
	db := lib.GetDatabase()
	var file File
	if err := db.Where(condition).First(&file).Error; err != nil {
		// Not found so returns a null pointer.
		return nil, err
	}

	return &file, nil
}
