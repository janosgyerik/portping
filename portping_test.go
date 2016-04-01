package portping

import (
	"testing"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	testHost = "localhost"
	testPort = "4269"
	knownNonexistentHost = "nonexistent.janosgyerik.com"
	defaultTimeout = 5 * time.Second
	testNetwork = "tcp"
)

func acceptN(t*testing.T, host, port string, count int) {
	ready := make(chan bool)
	go func() {
		ln, err := net.Listen(testNetwork, net.JoinHostPort(host, port))
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
	}()
	<-ready
}

func assertPingResult(t*testing.T, host, port string, expected bool, patterns ...string) {
	addr := net.JoinHostPort(host, port)
	err := Ping(testNetwork, addr, defaultTimeout)
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
		assertFormatResultContains(t, err, patterns...)
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
	addr := net.JoinHostPort(host, port)
	go PingN(testNetwork, addr, defaultTimeout, pingCount, c)

	failureCount := 0
	for i := 0; i < pingCount; i++ {
		err := <-c
		t.Logf("port ping %s [%d] -> %v", addr, i + 1, err)

		if err != nil {
			failureCount++
		}
	}

	successCount := pingCount - failureCount
	if expectedSuccessCount != successCount {
		t.Errorf("expected %d successful pings, but got only %d", expectedSuccessCount, successCount)
	}
}

func Test_ping_open_port(t*testing.T) {
	acceptN(t, testHost, testPort, 1)

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
	acceptN(t, testHost, testPort, pingCount)

	assertPingNSuccessCount(t, testHost, testPort, pingCount, pingCount)
}

func Test_ping5_all_fail(t*testing.T) {
	pingCount := 5
	successCount := 0
	assertPingNSuccessCount(t, testHost, testPort, pingCount, successCount)
}

func Test_ping5_partial_success(t*testing.T) {
	successCount := 3
	acceptN(t, testHost, testPort, successCount)

	pingCount := 5
	assertPingNSuccessCount(t, testHost, testPort, pingCount, successCount)
}

func assertFormatResultContains(t*testing.T, err error, patterns ...string) {
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

func pingAndAssertFormatResultContains(t*testing.T, host, port string, patterns ...string) {
	addr := net.JoinHostPort(host, port)
	assertFormatResultContains(t, Ping(testNetwork, addr, defaultTimeout), patterns...)
}

func Test_format_result_success(t*testing.T) {
	acceptN(t, testHost, testPort, 1)
	pingAndAssertFormatResultContains(t, testHost, testPort, "success")
}

func Test_format_result_connection_refused(t*testing.T) {
	pingAndAssertFormatResultContains(t, testHost, testPort, "connection refused")
}

func Test_format_result_invalid_port_m1(t*testing.T) {
	pingAndAssertFormatResultContains(t, testHost, "-1", "invalid port", "unknown port")
}

func Test_format_result_invalid_port_123456(t*testing.T) {
	pingAndAssertFormatResultContains(t, testHost, "123456", "invalid port", "unknown port")
}

func Test_format_result_nonexistent_host(t*testing.T) {
	host := knownNonexistentHost
	pingAndAssertFormatResultContains(t, host, testPort, fmt.Sprintf("lookup %s: no such host", host))
}
