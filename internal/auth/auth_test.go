package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPasswordAndCheck(t *testing.T) {
	password := "supersecret123"

	// Hash the password
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	// Correct password should match
	ok, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("error checking password hash: %v", err)
	}
	if !ok {
		t.Errorf("expected password to match hash, got false")
	}

	// Incorrect password should not match
	ok, err = CheckPasswordHash("wrongpassword", hash)
	if err != nil {
		t.Fatalf("error checking password hash: %v", err)
	}
	if ok {
		t.Errorf("expected password to not match hash, got true")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "mysecretkey"
	userID := uuid.New()
	expiration := 2 * time.Hour

	// Create JWT
	token, err := MakeJWT(userID, secret, expiration)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	// Validate JWT
	returnedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}
	if returnedID != userID {
		t.Errorf("expected userID %v, got %v", userID, returnedID)
	}

	// Test invalid secret
	_, err = ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Errorf("expected error validating with wrong secret, got nil")
	}

	// Test expired token
	expiredToken, _ := MakeJWT(userID, secret, -time.Hour) // negative duration = expired
	_, err = ValidateJWT(expiredToken, secret)
	if err == nil {
		t.Errorf("expected error validating expired token, got nil")
	}
}
