package keyrack

import (
  "os"
  "io/ioutil"
  "fmt"
  "sync"
  "encoding/json"
)

type Database struct {
  Version uint8
  Data *Secret
  top *Group
  mutex sync.Mutex
  counter uint64
}

func NewDatabase() (db *Database) {
  db = &Database{Version: 1, counter: 2}
  db.top = NewGroup(1, "Top", &db.counter)
  return
}

func LoadDatabase(filename, password string) (db *Database, err error) {
  var (
    f *os.File
    dbJSON, groupJSON []byte
  )

  f, err = os.Open(filename)
  if err != nil {
    return
  }
  defer f.Close()

  dbJSON, err = ioutil.ReadAll(f)
  if err != nil {
    return
  }
  db = new(Database)
  err = json.Unmarshal(dbJSON, db)
  if err != nil {
    return
  }
  groupJSON, err = db.Data.Message(password)
  if err != nil {
    return
  }
  db.top = new(Group)
  err = json.Unmarshal(groupJSON, db.top)

  return
}

func (db *Database) Top() *Group {
  return db.top
}

func (db *Database) Save(filename, password string) (err error) {
  db.mutex.Lock()
  defer db.mutex.Unlock()

  if db.Data != nil && !db.Data.IsPasswordValid(password) {
    err = fmt.Errorf("invalid password")
    return
  }

  /* Serialize group to JSON and encrypt */
  var groupJSON []byte
  groupJSON, err = json.Marshal(db.top)
  if err != nil {
    return
  }

  db.Data, err = NewSecret(groupJSON, password)
  if err != nil {
    return
  }

  /* Serialize database to JSON */
  var dbJSON []byte
  dbJSON, err = json.Marshal(db)
  if err != nil {
    return
  }

  /* Write to the file */
  var f *os.File
  f, err = os.Create(filename)
  if err != nil {
    return
  }
  defer f.Close()

  var n int
  n, err = f.Write(dbJSON)
  if err != nil {
    return
  }
  if n != len(dbJSON) {
    err = fmt.Errorf("couldn't write to file")
    return
  }

  return
}
