package keyrack

import "fmt"

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

  // First check to see if there is another login with the same site/username
  for _, login = range group.Logins {
    if login.Site == site && login.Username == username {
      err = fmt.Errorf("there is already a login with the same site/username")
      return
    }
  }

  login, err = NewLogin(site, username, password, master)
  if err == nil {
    group.Logins = append(group.Logins, login)
  }
  return
}

func (group *Group) AddGroup(name string) (err error) {
  var subgroup *Group

  for _, subgroup = range group.Groups {
    if subgroup.Name == name {
      err = fmt.Errorf("there is already a group with the same name")
      return
    }
  }

  subgroup = NewGroup(name)
  group.Groups = append(group.Groups, subgroup)
  return
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
