package keyrack

import "sync/atomic"

/* FIXME: Counter doesn't unserialize properly */

type Group struct {
  Id uint64
  Name string
  Logins LoginList
  Groups GroupList
  Counter *uint64
}

func NewGroup(id uint64, name string, counter *uint64) (group *Group) {
  group = &Group{Id: id, Name: name, Counter: counter}
  group.Logins = make(LoginList, 0)
  group.Groups = make(GroupList, 0)
  return
}

func (group *Group) AddLogin(site, username, password, master string) (err error) {
  var (
    id uint64
    login *Login
  )

  id = atomic.AddUint64(group.Counter, 1)
  login, err = NewLogin(id, site, username, password, master)
  if err == nil {
    group.Logins = append(group.Logins, login)
  }
  return
}

func (group *Group) AddGroup(name string) {
  subgroup := NewGroup(atomic.AddUint64(group.Counter, 1), name, group.Counter)
  group.Groups = append(group.Groups, subgroup)
}

type GroupList []*Group

func (list GroupList) Len() int {
  return len(list)
}

func (list GroupList) Less(i, j int) bool {
  if list[i].Name == list[j].Name {
    return list[i].Id < list[j].Id
  } else {
    return list[i].Name < list[j].Name
  }
}

func (list GroupList) Swap(i, j int) {
  list[i], list[j] = list[j], list[i]
}
