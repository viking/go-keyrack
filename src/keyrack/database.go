package keyrack

import (
  "os"
  "io"
  "io/ioutil"
  "fmt"
  "bytes"
  "sync"
  "sync/atomic"
  "encoding/json"
  "crypto/rand"
  "crypto/aes"
  "crypto/cipher"
  "crypto/sha256"
  "crypto/hmac"
  "code.google.com/p/go.crypto/scrypt"
)

const (
  magic = "keyrack"
  headerSize = len(magic) + 1 /*version*/ + 32 /*salt*/ + aes.BlockSize /*iv*/
)

type Database struct {
  Filename string
  Version uint8
  Top *Group
  key []byte
  salt []byte
  iv []byte
  mutex sync.RWMutex
  counter uint64
}

func NewDatabase(filename, password string) (db *Database, err error) {
  db = &Database{Filename: filename}
  err = db.init(password)
  if err != nil {
    return
  }
  return
}

func (db *Database) Save(password string) (err error) {
  /* Check password */
  var key []byte
  key, err = db.generateKey(password)
  if !bytes.Equal(db.key, key) {
    err = fmt.Errorf("invalid password")
    return
  }

  db.mutex.Lock()
  defer db.mutex.Unlock()

  /* Serialize group to JSON */
  var plaintext []byte
  plaintext, err = json.Marshal(db.Top)
  if err != nil {
    return
  }

  /* Pad plaintext if necessary */
  blen := len(plaintext) % aes.BlockSize
  if blen > 0 {
    plaintext = append(plaintext, make([]byte, aes.BlockSize - blen)...)
  }

  /* Encrypt the plaintext */
  var block cipher.Block
  block, err = aes.NewCipher(db.key)
  if err != nil {
    return
  }

  ciphertext := make([]byte, len(plaintext))
  mode := cipher.NewCBCEncrypter(block, db.iv)
  mode.CryptBlocks(ciphertext, plaintext)

  /* Write to the file */
  var f *os.File
  f, err = os.Create(db.Filename)
  if err != nil {
    return
  }
  defer f.Close()

  /* Write header */
  header := db.header()

  var n int
  n, err = f.Write(header)
  if err != nil {
    return
  }

  /* Write HMAC */
  mac := hmac.New(sha256.New, []byte(password))
  mac.Write(header)
  sum := mac.Sum(nil)
  n, err = f.Write(sum)
  if err != nil {
    return
  }
  if n != len(sum) {
    err = fmt.Errorf("couldn't write hmac")
  }

  /* Write ciphertext */
  n, err = f.Write(ciphertext)
  if err != nil {
    return
  }
  if n != len(ciphertext) {
    err = fmt.Errorf("couldn't write ciphertext")
    return
  }

  return
}

func (db *Database) init(password string) (err error) {
  /* Read version and salt */
  db.mutex.RLock()

  var (f *os.File; file bool)
  f, err = os.Open(db.Filename)
  if err == nil {
    /* Check for empty file */
    var fi os.FileInfo
    fi, err = f.Stat()
    if err != nil {
      return
    }
    if fi.Size() > 0 {
      file = true
      defer db.mutex.RUnlock()
      defer f.Close()
    } else {
      f.Close()
    }
  }

  if file {
    /* Read header */
    header := make([]byte, headerSize)
    var n int
    n, err = io.ReadFull(f, header)
    if err != nil {
      return
    }
    if n != len(header) {
      err = fmt.Errorf("couldn't read header")
      return
    }

    /* Check magic */
    if !bytes.Equal([]byte(magic), header[:len(magic)]) {
      err = fmt.Errorf("invalid file format")
      return
    }

    /* Read HMAC */
    sum := make([]byte, 32)
    n, err = io.ReadFull(f, sum)
    if err != nil {
      return
    }
    if n != len(sum) {
      err = fmt.Errorf("couldn't read hmac")
      return
    }

    /* Verify HMAC */
    mac := hmac.New(sha256.New, []byte(password))
    mac.Write(header)
    actualSum := mac.Sum(nil)
    if !hmac.Equal(actualSum, sum) {
      err = fmt.Errorf("invalid password")
      return
    }

    /* Version */
    header = header[len(magic):]
    db.Version = uint8(header[0])
    header = header[1:]

    /* Salt */
    db.salt = header[:32]
    header = header[32:]

    /* IV */
    db.iv = header[:aes.BlockSize]
  } else {
    db.mutex.RUnlock()

    /* Version */
    db.Version = 1

    /* Salt */
    db.salt = make([]byte, 32)
    var n int
    n, err = io.ReadFull(rand.Reader, db.salt)
    if err != nil {
      return
    }
    if n != len(db.salt) {
      err = fmt.Errorf("couldn't get salt")
      return
    }

    /* IV */
    db.iv = make([]byte, aes.BlockSize)
    n, err = io.ReadFull(rand.Reader, db.iv)
    if err != nil {
      return
    }
    if n != len(db.iv) {
      err = fmt.Errorf("couldn't get iv")
      return
    }
  }

  db.key, err = db.generateKey(password)
  if err != nil {
    return
  }

  if !file {
    db.Top = NewGroup(atomic.AddUint64(&db.counter, 1), "Top", &db.counter)
    return
  }

  /* Unencrypt data */
  var data []byte
  data, err = ioutil.ReadAll(f)
  if err != nil {
    return
  }
  if len(data) % aes.BlockSize != 0 {
    err = fmt.Errorf("invalid ciphertext")
    return
  }

  var block cipher.Block
  block, err = aes.NewCipher(db.key)
  if err != nil {
    return
  }

  mode := cipher.NewCBCDecrypter(block, db.iv)
  mode.CryptBlocks(data, data)

  /* Unpad message */
  for i := len(data)-1; i > 0 && data[i] == byte(0); i-- {
    data = data[:i]
  }

  /* Unmarshal JSON */
  err = json.Unmarshal(data, &db.Top)

  return
}

func (db *Database) generateKey(password string) (key []byte, err error) {
  key, err = scrypt.Key([]byte(password), db.salt, 16384, 8, 1, 32)
  return
}

func (db *Database) header() (data []byte) {
  data = make([]byte, 0, headerSize)

  /* Add magic */
  data = append(data, []byte(magic)...)

  /* Add version */
  data = append(data, byte(db.Version))

  /* Add salt */
  data = append(data, db.salt...)

  /* Add initialization vector */
  data = append(data, db.iv...)

  return
}
