package models

import (
	"github.com/jinzhu/gorm"
	"github.com/nireo/upfi/lib"
)

type File struct {
	gorm.Model
	Filename    string
	UUID        string
	Description string
	Size 		int64
	UserID 		uint
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
func (file *File) Delete() {
	db := lib.GetDatabase()
	db.Delete(&file)
}