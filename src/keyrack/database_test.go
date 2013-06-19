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

func TestNewDatabase_NonExistantFile(t *testing.T) {
  db, err := NewDatabase("foo.dat", "supersecret")
  if err != nil {
    t.Fatal(err)
  }
  if db.Version != 1 {
    t.Errorf("expected %v, got %v", 1, db.Version)
  }
}

func TestNewDatabase_EmptyFile(t *testing.T) {
  f := func (filename string) {
    db, err := NewDatabase(filename, "supersecret")
    if err != nil {
      t.Fatal(err)
    }
    if db.Version != 1 {
      t.Errorf("expected %v, got %v", 1, db.Version)
    }
  }
  tempfile(t, f)
}

func TestNewDatabase_LoadValidFile(t *testing.T) {
  f := func(filename string) {
    var (err error; db, db2 *Database)

    db, err = NewDatabase(filename, "supersecret")
    if err != nil {
      t.Error(err)
      return
    }
    db.Top.AddLogin("Twitter", "dude", "secret123")
    err = db.Save("supersecret")
    if err != nil {
      t.Error(err)
      return
    }

    db2, err = NewDatabase(filename, "supersecret")
    if err != nil {
      t.Error(err)
      return
    }
    if db.Version != db2.Version {
      t.Errorf("expected %v, got %v", db.Version, db2.Version)
    }
    if !reflect.DeepEqual(db.Top, db2.Top) {
      t.Errorf("expected %v, got %v", db.Top, db2.Top)
    }
  }
  tempfile(t, f)
}

func TestNewDatabase_WrongPassword(t *testing.T) {
  f := func(filename string) {
    var (err error; db *Database)

    db, err = NewDatabase(filename, "supersecret")
    if err != nil {
      t.Error(err)
      return
    }
    db.Top.AddLogin("Twitter", "dude", "secret123")
    err = db.Save("supersecret")
    if err != nil {
      t.Error(err)
      return
    }

    _, err = NewDatabase(filename, "wrong")
    if err == nil {
      t.Error("expected error, got none")
    } else if err.Error() != "invalid password" {
      t.Error("expected %v, got %v", "invalid password", err.Error())
    }
  }
  tempfile(t, f)
}

func TestDatabase_Save_WrongPassword(t *testing.T) {
  f := func(filename string) {
    var (err error; db *Database)

    db, err = NewDatabase(filename, "supersecret")
    if err != nil {
      t.Error(err)
      return
    }
    db.Top.AddLogin("Twitter", "dude", "secret123")
    err = db.Save("wrong")
    if err == nil {
      t.Error("expected error, got none")
    } else if err.Error() != "invalid password" {
      t.Errorf("expected %v, got %v", "invalid password", err.Error())
    }
  }
  tempfile(t, f)
}
