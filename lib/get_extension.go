package lib

import "strings"

func GetFileExtension(fileName string) string {
	splitted := strings.Split(fileName, ".")
	return splitted[len(splitted)-1]
}
