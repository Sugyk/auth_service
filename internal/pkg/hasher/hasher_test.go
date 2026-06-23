package hasher

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewPasswordHasher(t *testing.T) {
	h := NewPasswordHasher(bcrypt.MinCost)

	if h == nil {
		t.Fatal("expected non-nil hasher")
	}
	if h.cost != bcrypt.MinCost {
		t.Errorf("expected cost %d, got %d", bcrypt.MinCost, h.cost)
	}
}

func TestHashPassword(t *testing.T) {
	h := NewPasswordHasher(bcrypt.MinCost)

	t.Run("returns non-empty hash", func(t *testing.T) {
		hash, err := h.HashPassword("secret123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hash == "" {
			t.Error("expected non-empty hash")
		}
	})

	t.Run("same password produces different hashes (salt)", func(t *testing.T) {
		hash1, _ := h.HashPassword("secret123")
		hash2, _ := h.HashPassword("secret123")
		if hash1 == hash2 {
			t.Error("expected different hashes due to salting")
		}
	})

	t.Run("hash is not equal to original password", func(t *testing.T) {
		password := "secret123"
		hash, _ := h.HashPassword(password)
		if hash == password {
			t.Error("hash should not equal original password")
		}
	})

	t.Run("empty password is hashed without error", func(t *testing.T) {
		_, err := h.HashPassword("")
		if err != nil {
			t.Errorf("unexpected error for empty password: %v", err)
		}
	})

	t.Run("returns error on invalid cost", func(t *testing.T) {
		badHasher := NewPasswordHasher(100) // bcrypt max cost is 31
		_, err := badHasher.HashPassword("secret123")
		if err == nil {
			t.Error("expected error for invalid cost")
		}
	})
}

func TestCompareHashAndPassword(t *testing.T) {
	h := NewPasswordHasher(bcrypt.MinCost)

	t.Run("correct password matches hash", func(t *testing.T) {
		hash, _ := h.HashPassword("correct-password")
		if !h.CompareHashAndPassword("correct-password", hash) {
			t.Error("expected match for correct password")
		}
	})

	t.Run("wrong password does not match hash", func(t *testing.T) {
		hash, _ := h.HashPassword("correct-password")
		if h.CompareHashAndPassword("wrong-password", hash) {
			t.Error("expected no match for wrong password")
		}
	})

	t.Run("empty password does not match non-empty hash", func(t *testing.T) {
		hash, _ := h.HashPassword("secret123")
		if h.CompareHashAndPassword("", hash) {
			t.Error("expected no match for empty password")
		}
	})

	t.Run("invalid hash returns false", func(t *testing.T) {
		if h.CompareHashAndPassword("secret123", "not-a-valid-hash") {
			t.Error("expected false for invalid hash")
		}
	})
}
