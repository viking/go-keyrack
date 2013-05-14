package keyrack

import (
  "testing"
)

func TestGroupList_Len(t *testing.T) {
  list := &GroupList{[]*Group{NewGroup(1, "Foo")}}
  if list.Len() != 1 {
    t.Errorf("expected 1, got %d", list.Len())
  }
}

func TestGroupList_Less(t *testing.T) {
  list := &GroupList{[]*Group{
    NewGroup(1, "Foo"),
    NewGroup(2, "Bar"),
    NewGroup(3, "Bar"),
  }}
  if list.Less(0, 1) {
    t.Errorf("expected %v to be less than %v", list.Groups[1], list.Groups[0])
  }
  if !list.Less(1, 0) {
    t.Errorf("expected %v to be less than %v", list.Groups[1], list.Groups[0])
  }
  if !list.Less(1, 2) {
    t.Errorf("expected %v to be less than %v", list.Groups[1], list.Groups[2])
  }
  if list.Less(2, 1) {
    t.Errorf("expected %v to be less than %v", list.Groups[1], list.Groups[2])
  }
}

func TestGroupList_Swap(t *testing.T) {
  group_1 := NewGroup(1, "Foo")
  group_2 := NewGroup(2, "Bar")
  list := GroupList{[]*Group{group_1, group_2}}
  list.Swap(0, 1)
  if list.Groups[0] != group_2 || list.Groups[1] != group_1 {
    t.Error("the values weren't swapped")
  }
}
