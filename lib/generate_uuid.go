package lib

import uuid "github.com/satori/go.uuid"

func GenerateUUID() string {
	newV4, err := uuid.NewV4()

	if err != nil {
		panic(err)
	}

	return newV4.String()
}