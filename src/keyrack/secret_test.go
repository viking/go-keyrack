package keyrack

import (
  "testing"
  "bytes"
)

func TestSecret_Message(t *testing.T) {
  var (
    secret *Secret
    message, result []byte
    err error
  )

  message = []byte("secret")
  secret, err = NewSecret(message, "password")
  if err != nil {
    t.Error(err)
  }

  result, err = secret.Message("password")
  if err != nil {
    t.Error(err)
  }
  if !bytes.Equal(message, result) {
    t.Errorf("expected %v, got %v", message, result)
  }
}

func TestSecret_Message_WrongPassword(t *testing.T) {
  var (
    secret *Secret
    message []byte
    err error
  )

  message = []byte("secret")
  secret, err = NewSecret(message, "password")
  if err != nil {
    t.Error(err)
  }

  _, err = secret.Message("wrong")
  if err == nil {
    t.Error("expected an error, but there wasn't one")
  }
}
