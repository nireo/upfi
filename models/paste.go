package models

import (
	"github.com/nireo/upfi/lib"
	"gorm.io/gorm"
)

type Paste struct {
	gorm.Model
	UserID      uint
	Content     string
	Title       string
	Description string
	UUID        string
}

func (paste *Paste) Serialize() lib.JSON {
	return lib.JSON{
		"uuid":        paste.UUID,
		"description": paste.Description,
		"title":       paste.Title,
		"content":     paste.Content,
	}
}

func FindOnePaste(condition interface{}) (*Paste, error) {
	db := lib.GetDatabase()
	var paste Paste
	if err := db.Where(condition).First(paste).Error; err != nil {
		return &paste, err
	}

	return &paste, nil
}
