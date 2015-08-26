package main

import (
	"bytes"
	"github.com/viking/go-keyrack/keyrack"
	"testing"
)

func TestGroupView_Render_NoLogins(t *testing.T) {
	group := keyrack.NewGroup("blah")
	group.AddGroup("foo")
	group.AddGroup("bar")
	view := GroupView{group}
	result, err := view.Render()
	if err == nil {
		expected := []byte("Group blah\n\nG0. foo\nG1. bar\n")
		if !bytes.Equal(result, expected) {
			t.Errorf("expected <%+v>, got <%+v>", []byte(expected), []byte(result))
		}
	} else {
		t.Error(err)
	}
}

func TestGroupView_Render_NoGroups(t *testing.T) {
	group := keyrack.NewGroup("blah")
	group.AddLogin("foo", "bro", "secret")
	group.AddLogin("bar", "dude", "secret")
	view := GroupView{group}
	result, err := view.Render()
	if err == nil {
		expected := []byte("Group blah\n\nL0. foo (bro)\nL1. bar (dude)\n")
		if !bytes.Equal(result, expected) {
			t.Errorf("expected <%+v>, got <%+v>", []byte(expected), []byte(result))
		}
	} else {
		t.Error(err)
	}
}
