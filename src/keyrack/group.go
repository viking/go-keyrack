package keyrack

type Group struct {
  Id uint64
  Name string
  Logins LoginList
}

func NewGroup(id uint64, name string) (group *Group) {
  group = &Group{Id: id, Name: name}
  group.Logins = make(LoginList, 0)
  return
}

type GroupList []*Group

func (list GroupList) Len() int {
  return len([]*Group(list))
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
