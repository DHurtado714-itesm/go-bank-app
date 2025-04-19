package utils

import (
	"testing"
)

func TestHashPasswordAndCheck_Success(t *testing.T) {
	// Arrange
	password := "supersecure123"

	// Act
	hashedPassword, err := HashPassword(password)
	match := CheckPasswordHash(password, hashedPassword)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if hashedPassword == password {
		t.Errorf("expected hashed password to differ from original")
	}
	if !match {
		t.Errorf("expected password and hash to match, but did not")
	}
}

func TestCheckPasswordHash_FailsOnWrongPassword(t *testing.T) {
	// Arrange
	password := "supersecure123"
	wrongPassword := "totallyWrongPassword"

	hashedPassword, err := HashPassword(password)

	// Act
	match := CheckPasswordHash(wrongPassword, hashedPassword)

	// Assert
	if err != nil {
		t.Fatalf("expected no error hashing password, got: %v", err)
	}
	if match {
		t.Errorf("expected password and hash to not match, but they did")
	}
}
