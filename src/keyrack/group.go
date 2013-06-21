package keyrack

type Group struct {
  Name string
  Logins LoginList
  Groups GroupList
}

func NewGroup(name string) (group *Group) {
  group = &Group{Name: name}
  group.Logins = make(LoginList, 0)
  group.Groups = make(GroupList, 0)
  return
}

func (group *Group) AddLogin(site, username, password, master string) (err error) {
  var login *Login

  login, err = NewLogin(site, username, password, master)
  if err == nil {
    group.Logins = append(group.Logins, login)
  }
  return
}

func (group *Group) AddGroup(name string) {
  subgroup := NewGroup(name)
  group.Groups = append(group.Groups, subgroup)
}

type GroupList []*Group

func (list GroupList) Len() int {
  return len(list)
}

func (list GroupList) Less(i, j int) bool {
  return list[i].Name < list[j].Name
}

func (list GroupList) Swap(i, j int) {
  list[i], list[j] = list[j], list[i]
}
