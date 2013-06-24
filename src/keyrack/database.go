package keyrack

import (
  "os"
  "io"
  "io/ioutil"
  "fmt"
  "sync"
  "encoding/json"
  "crypto/rand"
)

type Database struct {
  Version uint8
  Data *Secret
  private struct {
    Top *Group
    Key []byte
  }
  mutex sync.Mutex
}

func NewDatabase() (db *Database, err error) {
  db = &Database{Version: 1}

  var n int
  db.private.Key = make([]byte, 32)
  n, err = io.ReadFull(rand.Reader, db.private.Key)
  if err != nil {
    return
  }
  if n != len(db.private.Key) {
    err = fmt.Errorf("couldn't generate key")
    return
  }

  db.private.Top = NewGroup("Top")
  return
}

func LoadDatabase(filename string, password []byte) (db *Database, err error) {
  var (
    f *os.File
    dbJSON, privateJSON []byte
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
  privateJSON, err = db.Data.Message(password)
  if err != nil {
    return
  }
  err = json.Unmarshal(privateJSON, &db.private)

  return
}

func (db *Database) Key() []byte {
  return db.private.Key
}

func (db *Database) Top() *Group {
  return db.private.Top
}

func (db *Database) Save(filename string, password []byte) (err error) {
  db.mutex.Lock()
  defer db.mutex.Unlock()

  if db.Data != nil && !db.Data.IsPasswordValid(password) {
    err = fmt.Errorf("invalid password")
    return
  }

  /* Encrypt all the logins */
  db.encryptLogins(db.private.Top)

  /* Serialize group to JSON and encrypt */
  var privateJSON []byte
  privateJSON, err = json.Marshal(db.private)
  if err != nil {
    return
  }

  db.Data, err = NewSecret(privateJSON, password)
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

func (db *Database) encryptLogins(group *Group) (err error) {
  for _, login := range group.Logins {
    err = login.Encrypt(db.private.Key)
    if err != nil {
      return
    }
  }
  for _, subgroup := range group.Groups {
    err = db.encryptLogins(subgroup)
    if err != nil {
      return
    }
  }
  return
}
