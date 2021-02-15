package utils

import "testing"

func TestGetEnv(t *testing.T) {
	GetEnv("APP_NAME")
}

func TestCleanPhoneNumber(t *testing.T) {
	cleanedPhoneNumber := CleanPhoneNumber("08199999")

	if cleanedPhoneNumber != "628199999" {
		t.Error("unexpected response")
	}
}
