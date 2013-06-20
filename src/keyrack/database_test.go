package keyrack

import (
  "testing"
  "io/ioutil"
  "os"
  "reflect"
)

func tempfile(t *testing.T, fn func (filename string)) {
  var (
    f *os.File
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
  /*defer os.Remove(f.Name())*/
  t.Log(f.Name())

  fn(f.Name())
}

func TestNewDatabase(t *testing.T) {
  db := NewDatabase()
  if db.Version != 1 {
    t.Errorf("expected %v, got %v", 1, db.Version)
  }
}

func TestLoadDatabase_EmptyFile(t *testing.T) {
  f := func (filename string) {
    _, err := LoadDatabase(filename, "supersecret")
    if err == nil {
      t.Error("expected an error, but there wasn't one")
    }
  }
  tempfile(t, f)
}

func TestLoadDatabase_LoadValidFile(t *testing.T) {
  f := func(filename string) {
    var (err error; db, db2 *Database)

    db = NewDatabase()
    top := db.Top()
    err = top.AddLogin("Twitter", "dude", "secret123", "supersecret")
    if err != nil {
      t.Error(err)
      return
    }
    err = db.Save(filename, "supersecret")
    if err != nil {
      t.Error(err)
      return
    }

    db2, err = LoadDatabase(filename, "supersecret")
    if err != nil {
      t.Error(err)
      return
    }
    if db.Version != db2.Version {
      t.Errorf("expected %v, got %v", db.Version, db2.Version)
    }
    top2 := db2.Top()
    if !reflect.DeepEqual(top, top2) {
      t.Errorf("expected %v, got %v", top, top2)
    }
  }
  tempfile(t, f)
}

func TestNewDatabase_WrongPassword(t *testing.T) {
  f := func(filename string) {
    var (err error; db *Database)

    db = NewDatabase()
    if err != nil {
      t.Error(err)
      return
    }
    top := db.Top()
    err = top.AddLogin("Twitter", "dude", "secret123", "supersecret")
    if err != nil {
      t.Error(err)
      return
    }
    err = db.Save(filename, "supersecret")
    if err != nil {
      t.Error(err)
      return
    }

    _, err = LoadDatabase(filename, "wrong")
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
    var (err error; db *Database)

    db = NewDatabase()
    if err != nil {
      t.Error(err)
      return
    }
    top := db.Top()
    err = top.AddLogin("Twitter", "dude", "secret123", "supersecret")
    if err != nil {
      t.Error(err)
      return
    }
    err = db.Save(filename, "supersecret")
    if err != nil {
      t.Error(err)
      return
    }
    err = db.Save(filename, "wrong")
    if err == nil {
      t.Error("expected error, got none")
    } else if err.Error() != "invalid password" {
      t.Errorf("expected %v, got %v", "invalid password", err.Error())
    }
  }
  tempfile(t, f)
}
