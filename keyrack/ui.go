package main

import (
	"fmt"
	"github.com/viking/keyrack"
	"math"
	"strconv"
	"strings"
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
		index := fmt.Sprintf("L%d", i+1)
		fmt.Printf("% *s. %s [%s]\n", width, index, login.Site, login.Username)
	}
	fmt.Println("Commands: login group up save quit")
}

// read input from user
func getInput(prompt string, echo bool) (input string, err error) {
	if len(prompt) > 0 {
		fmt.Printf("%s ", prompt)
	}
	if echo {
		_, err = fmt.Scanf("%s", &input)
	} else {
		input = string(GetPasswd())
	}
	return
}

// print password from login
func printPassword(password string) {
	fmt.Printf("%s\r", password)
	getch()
	fmt.Println(strings.Repeat(" ", len(password)))
}

// add login to group
func newLogin(group *keyrack.Group) (err error) {
	var site, username, password string

	site, err = getInput("Site:", true)
	if err != nil || site == "" {
		return
	}

	username, err = getInput("Username:", true)
	if err != nil || username == "" {
		return
	}

	password, err = getInput("Password:", false)
	if err != nil || password == "" {
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

// add subgroup to group
func newGroup(group *keyrack.Group) (err error) {
	var name string

	name, err = getInput("Name:", true)
	if err != nil || name == "" {
		return
	}

	err = group.AddGroup(name)
	if err != nil {
		if err == keyrack.ErrGroupExists {
			fmt.Println("Error:", err)
			err = nil
		}
	}
	return
}

// main menu
func menu(session *Session, group *keyrack.Group) (quit bool, err error) {
	for !quit && err == nil {
		printMenu(group)

		up := false
		for ok := false; !ok && err == nil; {
			var command string
			command, err = getInput("?", true)
			if err != nil {
				return
			}

			ok = true
			switch command {
			case "login":
				err = newLogin(group)

			case "group":
				err = newGroup(group)

			case "up":
				up = true

			case "save":
				var password string
				password, err = getInput("Password:", false)
				if err == nil {
					err = session.db.Save(session.filename, []byte(password))
					if err != nil {
						if err == keyrack.ErrInvalidPassword {
							fmt.Println("Error:", err)
							err = nil
						}
					}
				}

			case "quit":
				quit = true

			default:
				var i int
				switch command[0] {
				case 'G', 'g':
					numGroups := len(group.Groups)
					i, err = strconv.Atoi(command[1:])
					if ok = err == nil && numGroups > 0 && i > 0 && i <= numGroups; ok {
						subgroup := group.Groups[i-1]
						quit, err = menu(session, subgroup)
					}

				case 'L', 'l':
					numLogins := len(group.Logins)
					i, err = strconv.Atoi(command[1:])
					if ok = err == nil && numLogins > 0 && i > 0 && i <= numLogins; ok {
						login := group.Logins[i-1]
						err = session.db.DecryptLogin(login)
						if err == nil {
							printPassword(login.Password())
							login.Clear()
						}
					}

				default:
					ok = false
				}
			}
		}

		if up && session.db.Top() != group {
			break
		}
	}
	return
}
