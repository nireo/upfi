package models

import (
	"os"

	"github.com/jinzhu/gorm"
	"github.com/nireo/upfi/lib"
)

type File struct {
	gorm.Model
	Filename    string
	UUID        string
	Description string
	Size        int64
	UserID      uint
}

func (file *File) Serialize() lib.JSON {
	return lib.JSON{
		"filename":    file.Filename,
		"created_at":  file.CreatedAt,
		"description": file.Description,
		"uuid":        file.UUID,
	}
}

// add a method to delete since we need this in the html template
func (file *File) Delete(userId string) error {
	db := lib.GetDatabase()

	err := os.Remove("./files/" + userId + "/" + file.Filename)
	if err != nil {
		return err
	}

	db.Delete(&file)
	return nil
}

func FindOneFile(condition interface{}) (*File, error) {
	db := lib.GetDatabase()

	var file File
	if err := db.Where(condition).First(&file).Error; err != nil {
		return &file, err
	}

	return &file, nil
}
