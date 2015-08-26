package main

import (
	"fmt"
	"github.com/viking/go-keyrack/keyrack"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"regexp"
)

type TermIO struct {
	Input  io.Reader
	Output io.Writer
}

func (t TermIO) Read(p []byte) (n int, err error) {
	n, err = t.Input.Read(p)
	return
}

func (t TermIO) Write(p []byte) (n int, err error) {
	n, err = t.Output.Write(p)
	return
}

func main() {
	var (
		filename string
		site     string
		username string
		password string
		line     string
		buf      []byte
		db       *keyrack.Database
		oldState *terminal.State
		term     *terminal.Terminal
		termio   TermIO
		group    *keyrack.Group
		//trail     []*keyrack.Group
		groupView keyrack.GroupView
		matched   bool
		quit      bool
		err       error
	)

	if len(os.Args) != 2 {
		fmt.Printf("Syntax: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}
	filename = os.Args[1]

	// setup terminal
	oldState, err = terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	termio.Input = os.Stdin
	termio.Output = os.Stdout
	term = terminal.NewTerminal(termio, "> ")

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		db, err = keyrack.NewDatabase()
	} else {
		password, err = term.ReadPassword("Password: ")
		if err == nil {
			db, err = keyrack.LoadDatabase(filename, []byte(password))
		}
		password = ""
	}

	group = db.Top()
	for err == nil && !quit {
		groupView.Group = group
		buf, err = groupView.Render()
		if err != nil {
			break
		}

		_, err = term.Write(buf)
		if err != nil {
			break
		}

		line, err = term.ReadLine()
		if err != nil {
			break
		}

		switch line {
		case "q":
			password, err = term.ReadPassword("Password: ")
			if err == nil {
				err = db.Save(filename, []byte(password))
			}
			password = ""
			quit = true
		case "ng", "new group":
			term.SetPrompt("Group name: ")
			line, err = term.ReadLine()
			if err != nil {
				break
			}

			if line == "" {
				term.Write([]byte("Group creation cancelled.\n"))
			} else {
				err = group.AddGroup(line)
				if err != nil {
					if err == keyrack.ErrGroupExists {
						term.Write([]byte("Group already exists.\n"))
					} else {
						break
					}
				}
			}
			term.SetPrompt("> ")
		case "nl", "new login":
			term.SetPrompt("Site name: ")
			site, err = term.ReadLine()
			if err != nil {
				break
			}

			if site == "" {
				term.Write([]byte("Login creation cancelled.\n"))
			} else {
				term.SetPrompt("Username: ")
				username, err = term.ReadLine()
				if err != nil {
					break
				}

				if username == "" {
					term.Write([]byte("Login creation cancelled.\n"))
				} else {
					password, err = term.ReadPassword("Password: ")
					if err != nil {
						break
					}

					if password == "" {
						term.Write([]byte("Login creation cancelled.\n"))
					} else {
						err = group.AddLogin(site, username, password)
						password = ""
						if err != nil {
							if err == keyrack.ErrLoginExists {
								term.Write([]byte("Login already exists.\n"))
							} else {
								break
							}
						}
					}
				}
			}

			term.SetPrompt("> ")
		default:
			matched, err = regexp.MatchString("^G(\\d+)$", line)
			if err != nil {
				break
			}
			if matched {
				term.Write([]byte("Group!\n"))
				continue
			}

			matched, err = regexp.MatchString("^L(\\d+)$", line)
			if err != nil {
				break
			}
			if matched {
				term.Write([]byte("Group!\n"))
			}
		}
	}

	if err != nil {
		fmt.Println(err)
	}

	if quit {
		fmt.Println("Quitting...")
	}
}
