package main

import (
	"fmt"
	"testing"
	"tumblr/redis"
)

var (
	Iterations = 50000
)

func TestService(t *testing.T) {
	// println("dialing redis...")
	ports := [3]int{6379, 6380, 6381}

	conns := [len(ports)]*redis.Conn{}

	for i := 0; i < len(ports); i++ {
		addr := fmt.Sprintf("localhost:%d", ports[i])
		// fmt.Printf("   %s\n", addr)

		conn, err := redis.Dial(addr)

		if err != nil {
			t.Fatalf("Could not connect to %s. Error: %s", addr, err)
		}

		conns[i] = conn
	}
	// println("   ... picked up")

	// println("tested.")
}
