package main

import "log"

const host = ":8888"

func main() {
	s := NewServer(host)

	err := s.Serve()
	if err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
