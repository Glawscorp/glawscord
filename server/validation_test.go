package server

import (
	"testing"
)

func TestValidUsername(t *testing.T) {

	got := validUsername("user123_-")
	want := true

	if got != want {

		t.Errorf("got %t, want %t", got, want)
	}
}

// too short
func TestInvalidUsername_short(t *testing.T) {

	got := validUsername("u")
	want := false

	if got != want {

		t.Errorf("got %t, want %t", got, want)
	}
}

// too long
func TestInvalidUsername_long(t *testing.T) {

	got := validUsername("user_alskdjfoauuwegbonofijqpweu857t9q73502uerqiwjbdfpaudshp9qew8hpfqo8hej;oisdhfi0")
	want := false

	if got != want {

		t.Errorf("got %t, want %t", got, want)
	}
}

// invalid chars
func TestInvalidUsername_chars(t *testing.T) {

	got := validUsername("user`~")
	want := false

	if got != want {

		t.Errorf("got %t, want %t", got, want)
	}
}

// Password tests
func TestValidPassword(t *testing.T) {

	got := validPassword("Password123_")
	want := true

	if got != want {

		t.Errorf("got %t, want %t", got, want)
	}
}

// too short
func TestInvalidPassword_short(t *testing.T) {

	got := validPassword("P")
	want := false

	if got != want {

		t.Errorf("got %t, want %t", got, want)
	}
}

// too long
func TestInvalidPassword_long(t *testing.T) {

	got := validPassword("Password123456789101112131415161718192021222324252627282930a;lskdjf;awbegohweoih______________")
	want := false

	if got != want {

		t.Errorf("got %t, want %t", got, want)
	}
}

// invalid chars
func TestInvalidPassword_chars(t *testing.T) {

	got := validPassword("Password123`~(")
	want := false

	if got != want {

		t.Errorf("got %t, want %t", got, want)
	}
}
