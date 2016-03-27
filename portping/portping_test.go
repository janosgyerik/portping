package portping

import (
	"testing"
	"fmt"
	"net"
	"log"
)

const testHost = "localhost"

// TODO hopefully unused. Better ideas?
const testPort = 1234

const knownNonexistentHost = "nonexistent.janosgyerik.com"

func acceptN(host string, port int, count int) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	for i := 0; i < count; i++ {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		conn.Close()
	}
}

func ping(host string, port int) bool {
	err := Ping(host, port)
	fmt.Printf("%s:%d %v", host, port, err)
	return err == nil
}

func Test_ping_open_port(t*testing.T) {
	go acceptN(testHost, testPort, 1)

	if !ping(testHost, testPort) {
		t.Errorf("should be open")
	}

	// for sanity: acceptN should have shut down already
	if ping(testHost, testPort) {
		t.Errorf("should be closed")
	}
}

func Test_ping_unopen_port(t*testing.T) {
	if ping(testHost, testPort) {
		t.Errorf("should be closed")
	}
}

func Test_ping_nonexistent_host(t*testing.T) {
	if ping(knownNonexistentHost, testPort) {
		t.Errorf("should be closed")
	}
}

func Test_ping_negative_port(t*testing.T) {
	if ping(testHost, -1) {
		t.Errorf("should be closed")
	}
}

func Test_ping_too_high_port(t*testing.T) {
	if ping(testHost, 123456) {
		t.Errorf("should be closed")
	}
}
