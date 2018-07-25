package models

import "testing"

func TestUser_NewUser(t *testing.T) {
	// Act
	user := NewUser(1, "kerkerj", "engineer")

	// Assert
	if user == nil {
		t.Fatal("user should not be nil")
	}

	if user.ID != 1 {
		t.Fatal("user's ID should be 1")
	}

	if user.Name != "kerkerj" {
		t.Fatal("user's name should be kerkerj")
	}
}
