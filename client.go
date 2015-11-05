package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Message struct {
	Source User
	Sent   time.Time
	Text   string
}

type User struct {
	net.Conn

	Name string

	scanner *bufio.Scanner
}

func (u User) Send(msg string) error {
	fmt.Fprintf(u, msg)

	return nil
}

func (u User) Listen() <-chan Message {
	in := make(chan Message)

	go func() {
		defer close(in)
		for u.scanner.Scan() {
			msg := u.scanner.Text()
			in <- Message{
				Source: u,
				Sent:   time.Now(),
				Text:   msg,
			}
		}
	}()

	return in
}
