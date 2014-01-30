// +build linux darwin
// Credit: https://github.com/howeyc/gopass

package main

import (
	"syscall"

	"code.google.com/p/go.crypto/ssh/terminal"
)

func getch() byte {
	if oldState, err := terminal.MakeRaw(0); err != nil {
		panic(err)
	} else {
		defer terminal.Restore(0, oldState)
	}

	var buf [1]byte
	if n, err := syscall.Read(0, buf[:]); n == 0 || err != nil {
		panic(err)
	}
	return buf[0]
}
