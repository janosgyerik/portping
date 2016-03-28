package main

import (
	"testing"
	"fmt"
	"net"
	"log"
	"strings"
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

func assertPingResult(host string, port int, t*testing.T, expected bool, pattern string) {
	err := Ping(host, port)

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("port ping %s -> %v", addr, err)

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

func assertPingSuccess(host string, port int, t*testing.T) {
	assertPingResult(host, port, t, true, "")
}

func assertPingFailure(host string, port int, t*testing.T, pattern string) {
	assertPingResult(host, port, t, false, pattern)
}

func assertPingNSuccessCount(host string, port int, t*testing.T, pingCount int, expectedSuccessCount int) {
	c := make(chan error)
	go PingN(host, port, pingCount, c)

	addr := fmt.Sprintf("%s:%d", host, port)

	successCount := 0
	for i := 0; i < pingCount; i++ {
		err := <-c
		log.Printf("port ping %s [%d] -> %v", addr, i + 1, err)

		if err == nil {
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
	assertPingFailure(testHost, testPort, t, "connection refused")
}

func Test_ping_unopen_port(t*testing.T) {
	assertPingFailure(testHost, testPort, t, "connection refused")
}

func Test_ping_nonexistent_host(t*testing.T) {
	assertPingFailure(knownNonexistentHost, testPort, t, "no such host")
}

func Test_ping_negative_port(t*testing.T) {
	assertPingFailure(testHost, -1, t, "invalid port")
}

func Test_ping_too_high_port(t*testing.T) {
	assertPingFailure(testHost, 123456, t, "invalid port")
}

func Test_ping5_all_success(t*testing.T) {
	pingCount := 3
	go acceptN(testHost, testPort, pingCount)

	assertPingNSuccessCount(testHost, testPort, t, pingCount, pingCount)
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

func assertFormatResult(host string, port int, t*testing.T, expected string) {
	actual := FormatResult(Ping(host, port))
	if expected != actual {
		t.Errorf("expected '%s' but got '%s'", expected, actual)
	}
}

func Test_format_result_success(t*testing.T) {
	go acceptN(testHost, testPort, 1)
	assertFormatResult(testHost, testPort, t, "success")
}

func Test_format_result_connection_refused(t*testing.T) {
	assertFormatResult(testHost, testPort, t, "getsockopt: connection refused")
}

func Test_format_result_invalid_port_m1(t*testing.T) {
	port := -1
	assertFormatResult(testHost, port, t, fmt.Sprintf("invalid port %d", port))
}

func Test_format_result_invalid_port_123456(t*testing.T) {
	port := 123456
	assertFormatResult(testHost, port, t, fmt.Sprintf("invalid port %d", port))
}

func Test_format_result_nonexistent_host(t*testing.T) {
	host := knownNonexistentHost
	assertFormatResult(host, testPort, t, fmt.Sprintf("lookup %s: no such host", host))
}
