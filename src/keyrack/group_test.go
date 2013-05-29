package keyrack

import (
  "testing"
)

func TestGroupList_Len(t *testing.T) {
  list := GroupList{NewGroup(1, "Foo")}
  if list.Len() != 1 {
    t.Errorf("expected 1, got %d", list.Len())
  }
}

func TestGroupList_Less(t *testing.T) {
  list := GroupList{
    NewGroup(1, "Foo"),
    NewGroup(2, "Bar"),
    NewGroup(3, "Bar"),
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
  group_1 := NewGroup(1, "Foo")
  group_2 := NewGroup(2, "Bar")
  list := GroupList{group_1, group_2}
  list.Swap(0, 1)
  if list[0] != group_2 || list[1] != group_1 {
    t.Error("the values weren't swapped")
  }
}
