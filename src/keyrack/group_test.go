package keyrack

import (
  "testing"
)

func TestGroup_AddLogin(t *testing.T) {
  counter := uint64(1)
  group := NewGroup(1, "Foo", &counter)
  group.AddLogin("Twitter", "dude", "secret")
  if len(group.Logins) != 1 {
    t.Errorf("expected 1, got %d", len(group.Logins))
  } else {
    if counter != 2 {
      t.Errorf("expected %v, got %v", 2, counter)
    }

    login := group.Logins[0]
    if login.Id != 2 {
      t.Errorf("expected %v, got %v", 2, login.Id)
    }
    if login.Site != "Twitter" {
      t.Errorf("expected %v, got %v", "Twitter", login.Site)
    }
    if login.Username != "dude" {
      t.Errorf("expected %v, got %v", "dude", login.Username)
    }
    if login.Password != "secret" {
      t.Errorf("expected %v, got %v", "secret", login.Password)
    }
  }
}

func TestGroup_AddGroup(t *testing.T) {
  counter := uint64(1)
  group := NewGroup(1, "Foo", &counter)
  group.AddGroup("Bar")
  if len(group.Groups) != 1 {
    t.Errorf("expected 1, got %d", len(group.Groups))
  } else {
    if counter != 2 {
      t.Errorf("expected %v, got %v", 2, counter)
    }

    subgroup := group.Groups[0]
    if subgroup.Id != 2 {
      t.Errorf("expected %v, got %v", 2, subgroup.Id)
    }
    if subgroup.Name != "Bar" {
      t.Errorf("expected %v, got %v", "Bar", subgroup.Name)
    }
  }
}

func TestGroupList_Len(t *testing.T) {
  counter := uint64(1)
  list := GroupList{NewGroup(1, "Foo", &counter)}
  if list.Len() != 1 {
    t.Errorf("expected 1, got %d", list.Len())
  }
}

func TestGroupList_Less(t *testing.T) {
  counter := uint64(1)
  list := GroupList{
    NewGroup(1, "Foo", &counter),
    NewGroup(2, "Bar", &counter),
    NewGroup(3, "Bar", &counter),
  }
  if list.Less(0, 1) {
    t.Errorf("expected %v to be less than %v", list[1], list[0])
  }
  if !list.Less(1, 0) {
    t.Errorf("expected %v to be less than %v", list[1], list[0])
  }
  if !list.Less(1, 2) {
    t.Errorf("expected %v to be less than %v", list[1], list[2])
  }
  if list.Less(2, 1) {
    t.Errorf("expected %v to be less than %v", list[1], list[2])
  }
}

func TestGroupList_Swap(t *testing.T) {
  counter := uint64(1)
  group_1 := NewGroup(1, "Foo", &counter)
  group_2 := NewGroup(2, "Bar", &counter)
  list := GroupList{group_1, group_2}
  list.Swap(0, 1)
  if list[0] != group_2 || list[1] != group_1 {
    t.Error("the values weren't swapped")
  }
}
