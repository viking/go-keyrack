package main

import (
	"bytes"
	"github.com/viking/go-keyrack/keyrack"
	"text/template"
)

var GroupTemplate string = "Group {{.Name}}\n\n{{range .Groups}}G{{.Index}}. {{.Name}}\n{{end}}{{range .Logins}}L{{.Index}}. {{.Site}} ({{.Username}})\n{{end}}"

type GroupView struct {
	Group *keyrack.Group
}

func (gv GroupView) Name() string {
	return gv.Group.Name
}

type GroupViewGroup struct {
	Index int
	Name  string
}

func (gv GroupView) Groups() (result []GroupViewGroup) {
	for i, group := range gv.Group.Groups {
		result = append(result, GroupViewGroup{i, group.Name})
	}
	return
}

type GroupViewLogin struct {
	Index    int
	Site     string
	Username string
}

func (gv GroupView) Logins() (result []GroupViewLogin) {
	for i, login := range gv.Group.Logins {
		result = append(result, GroupViewLogin{i, login.Site, login.Username})
	}
	return
}

func (gv GroupView) Render() (result []byte, err error) {
	tmpl, err := template.New("group").Parse(GroupTemplate)
	if err != nil {
		return
	}

	var buf []byte
	writer := bytes.NewBuffer(buf)
	err = tmpl.Execute(writer, gv)
	if err != nil {
		return
	}

	result = writer.Bytes()
	return
}
