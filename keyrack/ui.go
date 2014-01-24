package main

import (
	"fmt"
	"github.com/howeyc/gopass"
	"keyrack"
	"math"
)

// return the character width of the maximum index
func indexWidth(group *keyrack.Group) int {
	var (
		numGroups int
		numLogins int
		max       int
	)

	numGroups = len(group.Groups)
	numLogins = len(group.Logins)
	if numGroups > numLogins {
		max = numGroups
	} else {
		max = numLogins
	}
	return int(math.Floor(math.Log10(float64(max))))
}

// print a group menu
func printMenu(group *keyrack.Group) {
	fmt.Printf("=== %s\n", group.Name)

	width := indexWidth(group)

	for i, subgroup := range group.Groups {
		index := fmt.Sprintf("G%d", i+1)
		fmt.Printf("% *s. %s\n", width, index, subgroup.Name)
	}
	for i, login := range group.Logins {
		index := fmt.Sprintf("S%d", i+1)
		fmt.Printf("% *s. %s [%s]\n", width, index, login.Site, login.Username)
	}
	fmt.Println("Commands: new save quit")
}

// read input from user
func getInput(prompt string) (command string, err error) {
	fmt.Printf("%s ", prompt)
	_, err = fmt.Scanf("%s", &command)
	return
}

// read password from user
func getPassword() []byte {
	fmt.Printf("Password: ")
	return gopass.GetPasswd()
}

// add login to group
func newLogin(group *keyrack.Group) (err error) {
	var site, username, password string

	site, err = getInput("Site:")
	if err != nil || site == "" {
		return
	}

	username, err = getInput("Username:")
	if err != nil || username == "" {
		return
	}

	password = string(getPassword())
	if password == "" {
		return
	}

	err = group.AddLogin(site, username, password)
	if err != nil {
		if err == keyrack.ErrLoginExists {
			fmt.Println("Error:", err)
			err = nil
		}
	}
	return
}

func menu(session *Session, group *keyrack.Group) (quit bool, err error) {
	for !quit && err == nil {
		printMenu(group)

		for ok := false; !ok; {
			var command string
			command, err = getInput("?")
			if err != nil {
				return
			}

			ok = true
			switch command {
			case "new":
				err = newLogin(group)

			case "save":
				password := getPassword()
				err = session.db.Save(session.filename, password)
				if err != nil {
					if err == keyrack.ErrInvalidPassword {
						fmt.Println("Error:", err)
						err = nil
					}
				}

			case "quit":
				quit = true

			default:
				ok = false
			}
		}
	}
	return
}
