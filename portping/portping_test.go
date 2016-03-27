package portping

import (
	"testing"
)

// TODO use localhost
const host = "192.168.1.10"
const localhost = "localhost"

func Test_ping_open_port(t*testing.T) {
	// TODO listen on a port to ensure it's open
	open_port := 22

	if !Ping(host, open_port) {
		t.Errorf("should be open")
	}
}

func Test_ping_unopen_port(t*testing.T) {
	// TODO choose some random port guaranteed closed
	unopen_port := 123

	if Ping(host, unopen_port) {
		t.Errorf("should be closed")
	}
}

func Test_ping_nonexistent_host(t*testing.T) {
	if Ping("nonexistent.janosgyerik.com", 80) {
		t.Errorf("should be closed")
	}
}

func Test_ping_invalid_port(t*testing.T) {
	if Ping(localhost, -1) {
		t.Errorf("should be closed")
	}
}
