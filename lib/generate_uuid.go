package lib

import uuid "github.com/satori/go.uuid"

func GenerateUUID() string {
	uuid, err := uuid.NewV4()

	if err != nil {
		panic(err)
	}

	return uuid.String()
}