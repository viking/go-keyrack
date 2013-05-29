package keyrack

type Login struct {
  Id uint64
  Site string
  Username string
  Password string
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
