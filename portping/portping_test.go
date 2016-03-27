package portping

import (
	"testing"
	"fmt"
)

// TODO use localhost
const host = "192.168.1.10"
const localhost = "localhost"

func ping(host string, port int) bool {
	err := Ping(host, port)
	fmt.Printf("%s:%d %v", host, port, err)
	return err == nil
}

func Test_ping_open_port(t*testing.T) {
	// TODO listen on a port to ensure it's open
	open_port := 22

	if !ping(host, open_port) {
		t.Errorf("should be open")
	}
}

func Test_ping_unopen_port(t*testing.T) {
	// TODO choose some random port guaranteed closed
	unopen_port := 123

	if ping(host, unopen_port) {
		t.Errorf("should be closed")
	}
}

func Test_ping_nonexistent_host(t*testing.T) {
	if ping("nonexistent.janosgyerik.com", 80) {
		t.Errorf("should be closed")
	}
}

func Test_ping_invalid_port(t*testing.T) {
	if ping(localhost, -1) {
		t.Errorf("should be closed")
	}
}
