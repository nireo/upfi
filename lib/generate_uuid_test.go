package lib

import "testing"

func TestUUIDCreation(t *testing.T) {
	uuid := GenerateUUID()
	if uuid == "" {
		t.Error("A generated uuid is empty.")
		return
	}
}
