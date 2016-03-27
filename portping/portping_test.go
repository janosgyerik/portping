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

func assertPingResult(host string, port int, t*testing.T, expected bool) {
	err := Ping(host, port)

	actual := err == nil

	if expected != actual {
		var openOrClosed string
		if expected {
			openOrClosed = "open"
		} else {
			openOrClosed = "closed"
		}
		t.Errorf("%s:%d should be %s", host, port, openOrClosed)
	}
}

func assertPingSuccess(host string, port int, t*testing.T) {
	assertPingResult(host, port, t, true)
}

func assertPingFailure(host string, port int, t*testing.T) {
	assertPingResult(host, port, t, false)
}

func assertPingNSuccessCount(host string, port int, t*testing.T, pingCount int, expectedSuccessCount int) {
	c := make(chan error)
	go PingN(host, port, pingCount, c)

	successCount := 0
	for i := 0; i < pingCount; i++ {
		if <-c == nil {
			successCount++
		}
	}

	if expectedSuccessCount != successCount {
		t.Errorf("expected %d successful pings, but got only %d", expectedSuccessCount, successCount)
	}
}

func Test_ping_open_port(t*testing.T) {
	go acceptN(testHost, testPort, 1)

	assertPingSuccess(testHost, testPort, t)

	// for sanity: acceptN should have shut down already
	assertPingFailure(testHost, testPort, t)
}

func Test_ping_unopen_port(t*testing.T) {
	assertPingFailure(testHost, testPort, t)
}

func Test_ping_nonexistent_host(t*testing.T) {
	assertPingFailure(knownNonexistentHost, testPort, t)
}

func Test_ping_negative_port(t*testing.T) {
	assertPingFailure(testHost, -1, t)
}

func Test_ping_too_high_port(t*testing.T) {
	assertPingFailure(testHost, 123456, t)
}

func Test_ping5_all_success(t*testing.T) {
	count := 3
	go acceptN(testHost, testPort, count)

	assertPingNSuccessCount(testHost, testPort, t, count, count)
}

func Test_ping5_all_fail(t*testing.T) {
	pingCount := 5
	successCount := 0
	assertPingNSuccessCount(testHost, testPort, t, pingCount, successCount)
}

func Test_ping5_partial_success(t*testing.T) {
	successCount := 3
	go acceptN(testHost, testPort, successCount)

	pingCount := 5
	assertPingNSuccessCount(testHost, testPort, t, pingCount, successCount)
}

