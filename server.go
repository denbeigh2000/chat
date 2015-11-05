package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	// Protects Guests
	sync.RWMutex

	// Host to run on
	Host string

	users map[string]User
}

func NewServer(host string) Server {
	return Server{
		Host:  host,
		users: make(map[string]User),
	}
}

func (s Server) initUser(c net.Conn) (User, error) {
	sc := bufio.NewScanner(c)

	taken := true

	var username string

	for taken {
		fmt.Fprintf(c, "Enter username: ")
		sc.Scan()
		username = sc.Text()

		s.Lock()
		_, taken = s.users[username]
		if taken {
			fmt.Fprintf(c, "Username already taken\n")
		}
		s.Unlock()
	}

	user := User{c, username, sc}

	return user, nil
}

func (s Server) handleConn(c net.Conn) {
	user, err := s.initUser(c)
	if err != nil {
		log.Printf("Error initialising user: %v\n", err)
		return
	}

	s.Lock()
	s.users[user.Name] = user
	s.Unlock()

	defer delete(s.users, user.Name)

	s.DeliverInfo(fmt.Sprintf("New user: %v", user.Name))
	defer s.DeliverInfo(fmt.Sprintf("%v left the room", user.Name))

	for msg := range user.Listen() {
		s.Deliver(msg)
	}
}

func (s Server) DeliverInfo(msg string) {
	for _, user := range s.users {
		fmt.Fprintf(user, "[%v] %v\n", time.Now(), msg)
	}
}

func (s Server) Deliver(m Message) {
	s.RLock()
	defer s.RUnlock()

	for _, user := range s.users {
		if user != m.Source {
			fmt.Fprintf(user, "[%v] %v: %v\n", m.Sent, m.Source.Name, m.Text)
		}
	}
}

func (s Server) Serve() error {
	ln, err := net.Listen("tcp", s.Host)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n", err)
			continue
		}

		go s.handleConn(conn)
	}

	return nil
}
