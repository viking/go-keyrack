package keyrack

import (
  "io"
  "fmt"
  "bytes"
  "crypto/rand"
  "crypto/aes"
  "crypto/hmac"
  "crypto/sha256"
  "crypto/cipher"
  "code.google.com/p/go.crypto/scrypt"
)

type Secret struct {
  Salt []byte
  IV []byte
  Sum []byte
  Ciphertext []byte
}

func NewSecret(message []byte, password []byte) (s *Secret, err error) {
  s = new(Secret)

  /* Generate salt */
  s.Salt = make([]byte, 32)
  var n int
  n, err = io.ReadFull(rand.Reader, s.Salt)
  if err != nil {
    return
  }
  if n != len(s.Salt) {
    err = fmt.Errorf("couldn't get salt")
    return
  }

  /* Generate IV */
  s.IV = make([]byte, aes.BlockSize)
  n, err = io.ReadFull(rand.Reader, s.IV)
  if err != nil {
    return
  }
  if n != len(s.IV) {
    err = fmt.Errorf("couldn't get iv")
    return
  }

  /* Calculate HMAC */
  mac := hmac.New(sha256.New, password)
  mac.Write(s.Salt)
  mac.Write(s.IV)
  s.Sum = mac.Sum(nil)

  /* Generate derived key */
  var key []byte
  key, err = s.generateKey(password)
  if err != nil {
    return
  }

  /* Pad message if necessary */
  blen := len(message) % aes.BlockSize
  if blen > 0 {
    message = append(message, make([]byte, aes.BlockSize - blen)...)
  }

  /* Encrypt the message */
  var block cipher.Block
  block, err = aes.NewCipher(key)
  if err != nil {
    return
  }

  s.Ciphertext = make([]byte, len(message))
  mode := cipher.NewCBCEncrypter(block, s.IV)
  mode.CryptBlocks(s.Ciphertext, message)

  return
}

func (s *Secret) IsPasswordValid(password []byte) bool {
  mac := hmac.New(sha256.New, password)
  mac.Write(s.Salt)
  mac.Write(s.IV)
  sum := mac.Sum(nil)
  return bytes.Equal(s.Sum, sum)
}

func (s *Secret) Message(password []byte) (message []byte, err error) {
  if !s.IsPasswordValid(password) {
    err = fmt.Errorf("invalid password")
    return
  }
  if len(s.Ciphertext) % aes.BlockSize != 0 {
    err = fmt.Errorf("invalid ciphertext")
    return
  }

  /* Generate derived key */
  var key []byte
  key, err = s.generateKey(password)
  if err != nil {
    return
  }

  /* Decrypt message */
  var block cipher.Block
  block, err = aes.NewCipher(key)
  if err != nil {
    return
  }

  message = make([]byte, len(s.Ciphertext))
  mode := cipher.NewCBCDecrypter(block, s.IV)
  mode.CryptBlocks(message, s.Ciphertext)

  /* Unpad message */
  for i := len(message) - 1; i > 0 && message[i] == byte(0); i-- {
    message = message[:i]
  }
  return
}

func (s *Secret) generateKey(password []byte) (key []byte, err error) {
  key, err = scrypt.Key(password, s.Salt, 16384, 8, 1, 32)
  return
}
