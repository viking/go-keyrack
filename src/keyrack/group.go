package keyrack

type Group struct {
  Id uint64
  Name string
  Logins *LoginList
}

func NewGroup(id uint64, name string) (group *Group) {
  group = &Group{Id: id, Name: name}
  group.Logins = NewLoginList()
  return
}

type GroupList struct {
  Groups []*Group
}

func (list *GroupList) Len() int {
  return len(list.Groups)
}

func (list *GroupList) Less(i, j int) bool {
  if list.Groups[i].Name == list.Groups[j].Name {
    return list.Groups[i].Id < list.Groups[j].Id
  } else {
    return list.Groups[i].Name < list.Groups[j].Name
  }
}

func (list *GroupList) Swap(i, j int) {
  list.Groups[i], list.Groups[j] = list.Groups[j], list.Groups[i]
}
