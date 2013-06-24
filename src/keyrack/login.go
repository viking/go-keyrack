package keyrack

import "fmt"

type Login struct {
  Site string
  Username string
  Data *Secret
  password string
}

func NewLogin(site, username, password string) *Login {
  return &Login{Site: site, Username: username, password: password}
}

func (login *Login) Password() string {
  return login.password
}

func (login *Login) SetPassword(password string) {
  login.password = password
  login.Data = nil
}

func (login *Login) Encrypt(key []byte) (err error) {
  if login.Data == nil {
    if login.password == "" {
      err = fmt.Errorf("can't encrypt empty password")
      return
    }
    login.Data, err = NewSecret([]byte(login.password), key)
    if err == nil {
      login.password = ""
    }
  }
  return
}

func (login *Login) Decrypt(key []byte) (err error) {
  if login.Data != nil {
    var message []byte
    message, err = login.Data.Message(key)
    if err == nil {
      login.password = string(message)
    }
  }
  return
}

type LoginList []*Login

func (list LoginList) Len() int {
  return len(list)
}

func (list LoginList) Less(i, j int) bool {
  if list[i].Site == list[j].Site {
    return list[i].Username < list[j].Username
  } else {
    return list[i].Site < list[j].Site
  }
}

func (list LoginList) Swap(i, j int) {
  list[i], list[j] = list[j], list[i]
}
