package keyrack

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func tempfile(t *testing.T, fn func(filename string)) {
	var (
		f   *os.File
		err error
	)

	/* Create temporary file */
	f, err = ioutil.TempFile("", "keyrack")
	if err != nil {
		t.Error(err)
		return
	}
	err = f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	fn(f.Name())

	if t.Failed() {
		t.Log(f.Name())
	} else {
		os.Remove(f.Name())
	}
}

func TestNewDatabase(t *testing.T) {
	db, err := NewDatabase()
	if err != nil {
		t.Fatal(err)
	}
	if db.Version != 1 {
		t.Errorf("expected %v, got %v", 1, db.Version)
	}
}

func TestLoadDatabase_EmptyFile(t *testing.T) {
	f := func(filename string) {
		_, err := LoadDatabase(filename, []byte("supersecret"))
		if err == nil {
			t.Error("expected an error, but there wasn't one")
		}
	}
	tempfile(t, f)
}

func TestLoadDatabase_LoadValidFile(t *testing.T) {
	f := func(filename string) {
		var (
			err     error
			db, db2 *Database
		)

		db, err = NewDatabase()
		if err != nil {
			t.Fatal(err)
		}
		top := db.Top()
		err = top.AddLogin("Twitter", "dude", "secret123")
		if err != nil {
			t.Error(err)
			return
		}
		password := []byte("supersecret")
		err = db.Save(filename, password)
		if err != nil {
			t.Error(err)
			return
		}

		db2, err = LoadDatabase(filename, password)
		if err != nil {
			t.Error(err)
			return
		}
		if db.Version != db2.Version {
			t.Errorf("expected %v, got %v", db.Version, db2.Version)
		}
		top2 := db2.Top()
		if !reflect.DeepEqual(top, top2) {
			t.Errorf("expected %+v, got %+v", top, top2)
		}

		login := top2.Logins[0]
		err = db.DecryptLogin(login)
		if err != nil {
			t.Fatal(err)
		}

		if login.Password() != "secret123" {
			t.Errorf("expected %v, got %v", "secret123", login.Password())
		}
	}
	tempfile(t, f)
}

func TestNewDatabase_WrongPassword(t *testing.T) {
	f := func(filename string) {
		var (
			err error
			db  *Database
		)

		db, err = NewDatabase()
		if err != nil {
			t.Error(err)
			return
		}
		top := db.Top()
		err = top.AddLogin("Twitter", "dude", "secret123")
		if err != nil {
			t.Error(err)
			return
		}
		err = db.Save(filename, []byte("supersecret"))
		if err != nil {
			t.Error(err)
			return
		}

		_, err = LoadDatabase(filename, []byte("wrong"))
		if err == nil {
			t.Error("expected error, got none")
		} else if err.Error() != "invalid password" {
			t.Errorf("expected %v, got %v", "invalid password", err.Error())
		}
	}
	tempfile(t, f)
}

func TestDatabase_Save_WrongPassword(t *testing.T) {
	f := func(filename string) {
		var (
			err error
			db  *Database
		)

		db, err = NewDatabase()
		if err != nil {
			t.Error(err)
			return
		}
		top := db.Top()
		err = top.AddLogin("Twitter", "dude", "secret123")
		if err != nil {
			t.Error(err)
			return
		}
		err = db.Save(filename, []byte("supersecret"))
		if err != nil {
			t.Error(err)
			return
		}
		err = db.Save(filename, []byte("wrong"))
		if err == nil {
			t.Error("expected error, got none")
		} else if err.Error() != "invalid password" {
			t.Errorf("expected %v, got %v", "invalid password", err.Error())
		}
	}
	tempfile(t, f)
}
