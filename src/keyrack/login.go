package keyrack

type Login struct {
  Site string
  Username string
  Password *Secret
}

func NewLogin(site, username, password, master string) (login *Login, err error) {
  login = &Login{Site: site, Username: username}
  login.Password, err = NewSecret([]byte(password), master)
  return
}

func (login *Login) PasswordString(master string) (password string, err error) {
  var message []byte
  message, err = login.Password.Message(master)
  if err == nil {
    password = string(message)
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
