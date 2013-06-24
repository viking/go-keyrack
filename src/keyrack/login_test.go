package keyrack

import (
  "testing"
)

func TestLogin_Encrypt(t *testing.T) {
  login := NewLogin("Twitter", "dude", "secret")
  err := login.Encrypt([]byte("foo"))
  if err != nil {
    t.Error(err)
  }
  if login.Password() != "" {
    t.Error("password wasn't cleared")
  }
  if login.Data == nil {
    t.Error("secret wasn't created")
  }
}

func TestLogin_Encrypt_Twice(t *testing.T) {
  login := NewLogin("Twitter", "dude", "secret")
  err := login.Encrypt([]byte("foo"))
  if err != nil {
    t.Error(err)
  }
  expected := login.Data
  err = login.Encrypt([]byte("foo"))
  if err != nil {
    t.Error(err)
  }
  if login.Data != expected {
    t.Error("secret changed when it shouldn't have")
  }
}

func TestLogin_Encrypt_Empty(t *testing.T) {
  login := NewLogin("Twitter", "dude", "")
  err := login.Encrypt([]byte("foo"))
  if err == nil {
    t.Error("expected an error, but there wasn't one")
  }
}

func TestLogin_Decrypt(t *testing.T) {
  login := NewLogin("Twitter", "dude", "secret")
  err := login.Encrypt([]byte("foo"))
  if err != nil {
    t.Error(err)
  }
  err = login.Decrypt([]byte("foo"))
  if err != nil {
    t.Error(err)
  }
  if login.Password() != "secret" {
    t.Errorf("expected %v, got %v", "secret", login.Password())
  }
}

func TestLogin_Decrypt_Twice(t *testing.T) {
  login := NewLogin("Twitter", "dude", "secret")
  err := login.Encrypt([]byte("foo"))
  if err != nil {
    t.Error(err)
  }
  err = login.Decrypt([]byte("foo"))
  if err != nil {
    t.Error(err)
  }
  err = login.Decrypt([]byte("foo"))
  if err != nil {
    t.Error(err)
  }
  if login.Password() != "secret" {
    t.Errorf("expected %v, got %v", "secret", login.Password())
  }
}

func TestLogin_SetPassword(t *testing.T) {
  login := NewLogin("Twitter", "dude", "secret")
  err := login.Encrypt([]byte("foo"))
  if err != nil {
    t.Error(err)
  }
  login.SetPassword("secretfoo")
  if login.Password() != "secretfoo" {
    t.Errorf("expected %v, got %v", "secretfoo", login.Password())
  }
  if login.Data != nil {
    t.Errorf("secret wasn't cleared")
  }
}

func TestLoginList_Len(t *testing.T) {
  var err error
  list := make(LoginList, 1)
  list[0] = NewLogin("Twitter", "dude", "secret")
  if err != nil {
    t.Fatal(err)
  }
  if list.Len() != 1 {
    t.Errorf("expected 1, got %d", list.Len())
  }
}

func TestLoginList_Less(t *testing.T) {
  list := make(LoginList, 4)
  list[0] = NewLogin("Twitter", "dude", "secret")
  list[1] = NewLogin("Facebook", "foo@example.com", "password")
  list[2] = NewLogin("Facebook", "bar@example.com", "password")

  if list.Less(0, 1) {
    t.Errorf("expected %v to be less than %v", list[1], list[0])
  }
  if !list.Less(1, 0) {
    t.Errorf("expected %v to be less than %v", list[1], list[0])
  }
  if list.Less(1, 2) {
    t.Errorf("expected %v to be less than %v", list[2], list[1])
  }
  if !list.Less(2, 1) {
    t.Errorf("expected %v to be less than %v", list[2], list[1])
  }
}

func TestLoginList_Swap(t *testing.T) {
  login_1 := NewLogin("Twitter", "dude", "secret")
  login_2 := NewLogin("Facebook", "foo@example.com", "password")

  list := LoginList{login_1, login_2}
  list.Swap(0, 1)
  if list[0] != login_2 || list[1] != login_1 {
    t.Error("the values weren't swapped")
  }
}
