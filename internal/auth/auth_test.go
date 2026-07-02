package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "my-secret"

	token, err := MakeJWT(userID, secret, time.Minute*5)
	if err != nil {
		t.Fatal(err)
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatal(err)
	}

	if gotID != userID {
		t.Errorf("expected %v, got %v", userID, gotID)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "my-secret"

	token, err := MakeJWT(userID, secret, -time.Minute) // já expirado
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "secret-one", time.Minute*5)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateJWT(token, "secret-two")
	if err == nil {
		t.Error("expected error for wrong secret, got nil")
	}
}
