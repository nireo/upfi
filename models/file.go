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
	UserID      uint
	Extension   string `json:"extension"`
	MIME        string `json:"mime"`
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

	err := os.Remove("./files/" + userID + "/" + file.Filename)
	if err != nil {
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
