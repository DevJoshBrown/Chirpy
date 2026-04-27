package auth

import "testing"

func TestHashPasswordSucceeds(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
}

func TestHashPasswordProducesDifferentHashes(t *testing.T) {
	// argon2id uses a random salt, so the same password should never produce the same hash twice.
	password := "hunter2"
	h1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("first hash failed: %v", err)
	}
	h2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("second hash failed: %v", err)
	}
	if h1 == h2 {
		t.Errorf("expected different hashes for same password, got identical: %q", h1)
	}
}

func TestCheckPasswordHashMatches(t *testing.T) {
	password := "s3cret!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}
	ok, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if !ok {
		t.Error("expected password to match its own hash")
	}
}

func TestCheckPasswordHashRejectsWrongPassword(t *testing.T) {
	hash, err := HashPassword("right-password")
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}
	ok, err := CheckPasswordHash("wrong-password", hash)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if ok {
		t.Error("expected wrong password to be rejected")
	}
}

func TestCheckPasswordHashRejectsMalformedHash(t *testing.T) {
	// A malformed hash should produce an error, not just `false`.
	_, err := CheckPasswordHash("anything", "not-a-real-hash")
	if err == nil {
		t.Error("expected error for malformed hash, got nil")
	}
}

// Table-driven version of the match/reject cases — same coverage, less repetition.
// Optional — keep the discrete tests above OR this one, not both.
func TestCheckPasswordHashTable(t *testing.T) {
	hash, err := HashPassword("the-real-password")
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}

	cases := []struct {
		name     string
		password string
		want     bool
	}{
		{"correct password matches", "the-real-password", true},
		{"wrong password rejected", "wrong-password", false},
		{"empty password rejected", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CheckPasswordHash(tc.password, hash)
			if err != nil {
				t.Fatalf("CheckPasswordHash error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
