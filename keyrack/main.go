package main

import (
	"fmt"
	"keyrack"
	"os"
)

func main() {
	var (
		filename string
		db       *keyrack.Database
		err      error
	)

	if len(os.Args) != 2 {
		fmt.Printf("Syntax: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}
	filename = os.Args[1]

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		db, err = keyrack.NewDatabase()
	} else {
		password := getPassword()
		db, err = keyrack.LoadDatabase(os.Args[1], password)

		// scrub password
		for i := range password {
			password[i] = 0
		}
	}

	if err == nil {
		ui := &UI{db, filename}
		_, err = ui.menu(db.Top())
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
