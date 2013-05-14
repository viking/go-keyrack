package keyrack

import (
  "testing"
)

func TestLoginList_Len(t *testing.T) {
  list := &LoginList{[]*Login{{1, "Twitter", "dude", "secret"}}}
  if list.Len() != 1 {
    t.Errorf("expected 1, got %d", list.Len())
  }
}

func TestLoginList_Less(t *testing.T) {
  list := &LoginList{[]*Login{
    {1, "Twitter", "dude", "secret"},
    {2, "Facebook", "foo@example.com", "password"},
    {3, "Facebook", "bar@example.com", "password"},
    {4, "Facebook", "bar@example.com", "password"},
  }}
  if list.Less(0, 1) {
    t.Errorf("expected %v to be less than %v", list.Logins[1], list.Logins[0])
  }
  if !list.Less(1, 0) {
    t.Errorf("expected %v to be less than %v", list.Logins[1], list.Logins[0])
  }
  if list.Less(1, 2) {
    t.Errorf("expected %v to be less than %v", list.Logins[2], list.Logins[1])
  }
  if !list.Less(2, 1) {
    t.Errorf("expected %v to be less than %v", list.Logins[2], list.Logins[1])
  }
  if !list.Less(2, 3) {
    t.Errorf("expected %v to be less than %v", list.Logins[2], list.Logins[3])
  }
  if list.Less(3, 2) {
    t.Errorf("expected %v to be less than %v", list.Logins[2], list.Logins[3])
  }
}

func TestLoginList_Swap(t *testing.T) {
  login_1 := &Login{1, "Twitter", "dude", "secret"}
  login_2 := &Login{2, "Facebook", "foo@example.com", "password"}
  list := &LoginList{[]*Login{login_1, login_2}}
  list.Swap(0, 1)
  if list.Logins[0] != login_2 || list.Logins[1] != login_1 {
    t.Error("the values weren't swapped")
  }
}
