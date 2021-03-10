package lib

import (
	uuid "github.com/satori/go.uuid"
)

// GenerateUUID uses the 'github.com/satori/go.uuid' package to construct a new V4 unique id.
func GenerateUUID() string {
	newV4 := uuid.NewV4()
	// Return the uuid as a string, so that it's easier to store into the database.
	return newV4.String()
}
