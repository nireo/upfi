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
}

func (file *File) Serialize() lib.JSON {
	return lib.JSON{
		"filename":    file.Filename,
		"created_at":  file.CreatedAt,
		"description": file.Description,
		"uuid":        file.UUID,
	}
}
