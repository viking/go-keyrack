package keyrack

type Login struct {
  Id uint64
  Site string
  Username string
  Password *Secret
}

func NewLogin(id uint64, site, username, password, master string) (login *Login, err error) {
  login = &Login{Id: id, Site: site, Username: username}
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
    if list[i].Username == list[j].Username {
      return list[i].Id < list[j].Id
    } else {
      return list[i].Username < list[j].Username
    }
  } else {
    return list[i].Site < list[j].Site
  }
}

func (list LoginList) Swap(i, j int) {
  list[i], list[j] = list[j], list[i]
}
