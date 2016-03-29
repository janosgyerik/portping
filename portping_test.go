package main

import (
	"testing"
	"fmt"
	"net"
	"strings"
)

const testHost = "localhost"

// TODO hopefully unused. Better ideas?
const testPort = 1234

const knownNonexistentHost = "nonexistent.janosgyerik.com"

func acceptN(t*testing.T, host string, port int, count int) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	for i := 0; i < count; i++ {
		conn, err := ln.Accept()
		if err != nil {
			t.Fatal(err)
		}
		conn.Close()
	}
}

func assertPingResult(t*testing.T, host string, port int, expected bool, pattern string) {
	err := Ping(host, port)

	addr := fmt.Sprintf("%s:%d", host, port)
	t.Logf("port ping %s -> %v", addr, err)

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

	if pattern != "" {
		errstr := err.Error()
		if !strings.Contains(errstr, pattern) {
			t.Errorf("the result was expected to contain %s, but was: %s", pattern, errstr)
		}
	}
}

func assertPingSuccess(t*testing.T, host string, port int) {
	assertPingResult(t, host, port, true, "")
}

func assertPingFailure(t*testing.T, host string, port int, pattern string) {
	assertPingResult(t, host, port, false, pattern)
}

func assertPingNSuccessCount(t*testing.T, host string, port int, pingCount int, expectedSuccessCount int) {
	c := make(chan error)
	go PingN(host, port, pingCount, c)

	addr := fmt.Sprintf("%s:%d", host, port)

	successCount := 0
	for i := 0; i < pingCount; i++ {
		err := <-c
		t.Logf("port ping %s [%d] -> %v", addr, i + 1, err)

		if err == nil {
			successCount++
		}
	}

	if expectedSuccessCount != successCount {
		t.Errorf("expected %d successful pings, but got only %d", expectedSuccessCount, successCount)
	}
}

func Test_ping_open_port(t*testing.T) {
	go acceptN(t, testHost, testPort, 1)

	assertPingSuccess(t, testHost, testPort)

	// for sanity: acceptN should have shut down already
	assertPingFailure(t, testHost, testPort, "connection refused")
}

func Test_ping_unopen_port(t*testing.T) {
	assertPingFailure(t, testHost, testPort, "connection refused")
}

func Test_ping_nonexistent_host(t*testing.T) {
	assertPingFailure(t, knownNonexistentHost, testPort, "no such host")
}

func Test_ping_negative_port(t*testing.T) {
	assertPingFailure(t, testHost, -1, "invalid port")
}

func Test_ping_too_high_port(t*testing.T) {
	assertPingFailure(t, testHost, 123456, "invalid port")
}

func Test_ping5_all_success(t*testing.T) {
	pingCount := 3
	go acceptN(t, testHost, testPort, pingCount)

	assertPingNSuccessCount(t, testHost, testPort, pingCount, pingCount)
}

func Test_ping5_all_fail(t*testing.T) {
	pingCount := 5
	successCount := 0
	assertPingNSuccessCount(t, testHost, testPort, pingCount, successCount)
}

func Test_ping5_partial_success(t*testing.T) {
	successCount := 3
	go acceptN(t, testHost, testPort, successCount)

	pingCount := 5
	assertPingNSuccessCount(t, testHost, testPort, pingCount, successCount)
}

func assertFormatResult(t*testing.T, host string, port int, expected string) {
	actual := FormatResult(Ping(host, port))
	if expected != actual {
		t.Errorf("expected '%s' but got '%s'", expected, actual)
	}
}

func Test_format_result_success(t*testing.T) {
	go acceptN(t, testHost, testPort, 1)
	assertFormatResult(t, testHost, testPort, "success")
}

func Test_format_result_connection_refused(t*testing.T) {
	assertFormatResult(t, testHost, testPort, "getsockopt: connection refused")
}

func Test_format_result_invalid_port_m1(t*testing.T) {
	port := -1
	assertFormatResult(t, testHost, port, fmt.Sprintf("invalid port %d", port))
}

func Test_format_result_invalid_port_123456(t*testing.T) {
	port := 123456
	assertFormatResult(t, testHost, port, fmt.Sprintf("invalid port %d", port))
}

func Test_format_result_nonexistent_host(t*testing.T) {
	host := knownNonexistentHost
	assertFormatResult(t, host, testPort, fmt.Sprintf("lookup %s: no such host", host))
}
