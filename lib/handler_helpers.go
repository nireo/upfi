package lib

import (
	"errors"

	"github.com/nireo/upfi/models"
)

// FindFileAndCheckOwnership is just a shortened version for a very common piece of code found in
// most of the file related handlers. And I thought this would make it more clear
func FindFileAndCheckOwnership(userID uint, fileID string) (*models.File, error) {
	file, err := models.FindOneFile(&models.File{UUID: fileID})
	if err != nil {
		return nil, err
	}

	if file.UserID != userID {
		return nil, errors.New("the user does not own this file.")
	}

	return file, nil
}
