package keyrack

import (
  "testing"
)

func TestLogin_PasswordString(t *testing.T) {
  var (
    login *Login
    password string
    err error
  )

  login, err = NewLogin("Twitter", "dude", "secret", "foo")
  if err != nil {
    t.Fatal(err)
  }

  password, err = login.PasswordString("foo")
  if err != nil {
    t.Error(err)
  }
  if password != "secret" {
    t.Errorf("expected %v, got %v", "secret", password)
  }
}

func TestLoginList_Len(t *testing.T) {
  var err error
  list := make(LoginList, 1)
  list[0], err = NewLogin("Twitter", "dude", "secret", "foo")
  if err != nil {
    t.Fatal(err)
  }
  if list.Len() != 1 {
    t.Errorf("expected 1, got %d", list.Len())
  }
}

func TestLoginList_Less(t *testing.T) {
  var err error
  list := make(LoginList, 4)
  list[0], err = NewLogin("Twitter", "dude", "secret", "foo")
  if err == nil {
    list[1], err = NewLogin("Facebook", "foo@example.com", "password", "foo")
  }
  if err == nil {
    list[2], err = NewLogin("Facebook", "bar@example.com", "password", "foo")
  }
  if err != nil {
    t.Fatal(err)
  }

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
  var (
    login_1, login_2 *Login
    err error
  )
  login_1, err = NewLogin("Twitter", "dude", "secret", "foo")
  if err == nil {
    login_2, err = NewLogin("Facebook", "foo@example.com", "password", "foo")
  }
  if err != nil {
    t.Fatal(err)
  }

  list := LoginList{login_1, login_2}
  list.Swap(0, 1)
  if list[0] != login_2 || list[1] != login_1 {
    t.Error("the values weren't swapped")
  }
}
