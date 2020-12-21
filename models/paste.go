package models

import (
	"github.com/nireo/upfi/lib"
	"gorm.io/gorm"
)

// Paste is a database struct, which also holds all the properties of gorm.Model
type Paste struct {
	gorm.Model
	UserID      uint
	Content     string
	Title       string
	Description string
	UUID        string
	Private     bool // If the private flag is provided, the paste is encrypted and also
}

// Serialize serializes a paste's data into json format
func (paste *Paste) Serialize() lib.JSON {
	return lib.JSON{
		"uuid":        paste.UUID,
		"description": paste.Description,
		"title":       paste.Title,
		"content":     paste.Content,
	}
}

// FindOnePaste takes a condition argument and returns a pointer to a paste,
// if it finds one matching the condition.
func FindOnePaste(condition interface{}) (*Paste, error) {
	db := lib.GetDatabase()
	var paste Paste
	if err := db.Where(condition).First(paste).Error; err != nil {
		// Paste not found so return a null pointer.
		return nil, err
	}

	return &paste, nil
}
