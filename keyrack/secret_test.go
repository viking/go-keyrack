package keyrack

import (
	"bytes"
	"testing"
)

func TestSecret_Message(t *testing.T) {
	var (
		secret          *Secret
		message, result []byte
		err             error
	)

	message = []byte("secret")
	secret, err = NewSecret(message, []byte("password"))
	if err != nil {
		t.Fatal(err)
	}

	result, err = secret.Message([]byte("password"))
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(message, result) {
		t.Errorf("expected %v, got %v", message, result)
	}
}

func TestSecret_Message_WrongPassword(t *testing.T) {
	var (
		secret  *Secret
		message []byte
		err     error
	)

	message = []byte("secret")
	secret, err = NewSecret(message, []byte("password"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = secret.Message([]byte("wrong"))
	if err == nil {
		t.Error("expected an error, but there wasn't one")
	}
}

func TestSecret_IsPasswordValid(t *testing.T) {
	secret, err := NewSecret([]byte("secret"), []byte("password"))
	if err != nil {
		t.Fatal(err)
	}

	if !secret.IsPasswordValid([]byte("password")) {
		t.Error("expected password to be valid, but wasn't")
	}
	if secret.IsPasswordValid([]byte("foobar")) {
		t.Error("expected password to be invalid, but wasn't")
	}
}
