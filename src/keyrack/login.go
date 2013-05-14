package keyrack

type Login struct {
  Id uint64
  Site string
  Username string
  Password string
}

type LoginList struct {
  Logins []*Login
}

func NewLoginList() (list *LoginList) {
  list = new(LoginList)
  list.Logins = make([]*Login, 0)
  return
}

func (list *LoginList) Len() int {
  return len(list.Logins)
}

func (list *LoginList) Less(i, j int) bool {
  if list.Logins[i].Site == list.Logins[j].Site {
    if list.Logins[i].Username == list.Logins[j].Username {
      return list.Logins[i].Id < list.Logins[j].Id
    } else {
      return list.Logins[i].Username < list.Logins[j].Username
    }
  } else {
    return list.Logins[i].Site < list.Logins[j].Site
  }
}

func (list *LoginList) Swap(i, j int) {
  list.Logins[i], list.Logins[j] = list.Logins[j], list.Logins[i]
}
