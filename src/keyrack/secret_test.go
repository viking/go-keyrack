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
