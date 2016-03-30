package main

import (
	"testing"
	"fmt"
	"net"
	"strings"
)

const testHost = "localhost"

// TODO hopefully unused. Better ideas?
const testPort = "1234"

const knownNonexistentHost = "nonexistent.janosgyerik.com"

func acceptN(t*testing.T, host, port string, count int, ready chan bool) {
	ln, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	ready <- true

	for i := 0; i < count; i++ {
		conn, err := ln.Accept()
		if err != nil {
			t.Fatal(err)
		}
		conn.Close()
	}
}

func assertPingResult(t*testing.T, host, port string, expected bool, patterns ...string) {
	err := Ping(host, port)

	addr := net.JoinHostPort(host, port)
	t.Logf("port ping %s -> %v", addr, err)

	actual := err == nil

	if actual != expected {
		var openOrClosed string
		if expected {
			openOrClosed = "open"
		} else {
			openOrClosed = "closed"
		}
		t.Errorf("%s should be %s", addr, openOrClosed)
	}

	if err != nil {
		result := FormatResult(err)
		foundMatch := false
		for _, pattern := range patterns {
			if strings.Contains(result, pattern) {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			t.Errorf("got '%s'; expected to contain one of '%s'", result, patterns)
		}
	}
}

func assertPingSuccess(t*testing.T, host, port string) {
	assertPingResult(t, host, port, true, "")
}

func assertPingFailure(t*testing.T, host, port string, patterns ...string) {
	assertPingResult(t, host, port, false, patterns...)
}

func assertPingNSuccessCount(t*testing.T, host, port string, pingCount int, expectedSuccessCount int) {
	c := make(chan error)
	go PingN(host, port, pingCount, c)

	addr := net.JoinHostPort(host, port)

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
	ready := make(chan bool)
	go acceptN(t, testHost, testPort, 1, ready)
	<-ready

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
	assertPingFailure(t, testHost, "-1", "invalid port", "unknown port")
}

func Test_ping_too_high_port(t*testing.T) {
	assertPingFailure(t, testHost, "123456", "invalid port", "unknown port")
}

func Test_ping5_all_success(t*testing.T) {
	pingCount := 3
	ready := make(chan bool)
	go acceptN(t, testHost, testPort, pingCount, ready)
	<-ready

	assertPingNSuccessCount(t, testHost, testPort, pingCount, pingCount)
}

func Test_ping5_all_fail(t*testing.T) {
	pingCount := 5
	successCount := 0
	assertPingNSuccessCount(t, testHost, testPort, pingCount, successCount)
}

func Test_ping5_partial_success(t*testing.T) {
	successCount := 3
	ready := make(chan bool)
	go acceptN(t, testHost, testPort, successCount, ready)
	<-ready

	pingCount := 5
	assertPingNSuccessCount(t, testHost, testPort, pingCount, successCount)
}

func assertFormatResultContains(t*testing.T, host, port string, patterns ...string) {
	result := FormatResult(Ping(host, port))
	foundMatch := false
	for _, pattern := range patterns {
		if strings.Contains(result, pattern) {
			foundMatch = true
			break
		}
	}
	if !foundMatch {
		t.Errorf("got '%s'; expected to contain one of '%s'", result, patterns)
	}
}

func Test_format_result_success(t*testing.T) {
	ready := make(chan bool)
	go acceptN(t, testHost, testPort, 1, ready)
	<-ready
	assertFormatResultContains(t, testHost, testPort, "success")
}

func Test_format_result_connection_refused(t*testing.T) {
	assertFormatResultContains(t, testHost, testPort, "connection refused")
}

func Test_format_result_invalid_port_m1(t*testing.T) {
	assertFormatResultContains(t, testHost, "-1", "invalid port", "unknown port")
}

func Test_format_result_invalid_port_123456(t*testing.T) {
	assertFormatResultContains(t, testHost, "123456", "invalid port", "unknown port")
}

func Test_format_result_nonexistent_host(t*testing.T) {
	host := knownNonexistentHost
	assertFormatResultContains(t, host, testPort, fmt.Sprintf("lookup %s: no such host", host))
}
