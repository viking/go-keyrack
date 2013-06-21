package keyrack

import (
  "testing"
)

func TestGroup_AddLogin(t *testing.T) {
  var (
    group *Group
    login *Login
    password string
    err error
  )

  group = NewGroup("Foo")
  err = group.AddLogin("Twitter", "dude", "secret", "foo")
  if err != nil {
    t.Fatal(err)
  }

  if len(group.Logins) != 1 {
    t.Errorf("expected 1, got %d", len(group.Logins))
  } else {
    login = group.Logins[0]
    if login.Site != "Twitter" {
      t.Errorf("expected %v, got %v", "Twitter", login.Site)
    }
    if login.Username != "dude" {
      t.Errorf("expected %v, got %v", "dude", login.Username)
    }
    password, err = login.PasswordString("foo")
    if err != nil {
      t.Error(err)
    }
    if password != "secret" {
      t.Errorf("expected %v, got %v", "secret", password)
    }
  }
}

func TestGroup_AddGroup(t *testing.T) {
  group := NewGroup("Foo")
  group.AddGroup("Bar")
  if len(group.Groups) != 1 {
    t.Errorf("expected 1, got %d", len(group.Groups))
  } else {
    subgroup := group.Groups[0]
    if subgroup.Name != "Bar" {
      t.Errorf("expected %v, got %v", "Bar", subgroup.Name)
    }
  }
}

func TestGroupList_Len(t *testing.T) {
  list := GroupList{NewGroup("Foo")}
  if list.Len() != 1 {
    t.Errorf("expected 1, got %d", list.Len())
  }
}

func TestGroupList_Less(t *testing.T) {
  list := GroupList{NewGroup("Foo"), NewGroup("Bar")}
  if list.Less(0, 1) {
    t.Errorf("expected %v to be less than %v", list[1], list[0])
  }
  if !list.Less(1, 0) {
    t.Errorf("expected %v to be less than %v", list[1], list[0])
  }
}

func TestGroupList_Swap(t *testing.T) {
  group_1 := NewGroup("Foo")
  group_2 := NewGroup("Bar")
  list := GroupList{group_1, group_2}
  list.Swap(0, 1)
  if list[0] != group_2 || list[1] != group_1 {
    t.Error("the values weren't swapped")
  }
}
