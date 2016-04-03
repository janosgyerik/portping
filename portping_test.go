package portping

import (
	"testing"
	"net"
	"strings"
	"time"
	"strconv"
)

const (
	testHost = "localhost"
	knownNonexistentHost = "nonexistent.janosgyerik.com"
	defaultTimeout = 5 * time.Second
	testNetwork = "tcp"
)

var testPort = findKnownAvailablePort()

func findKnownAvailablePort() string {
	tcpa, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		panic(err)
	}
	ln, err := net.ListenTCP("tcp", tcpa)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	local, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		panic("Unable to convert Addr to TCPAddr")
	}

	return strconv.Itoa(local.Port)
}

func acceptN(t*testing.T, host, port string, count int, ready chan <- bool, done chan <- bool) {
	ln, err := net.Listen(testNetwork, net.JoinHostPort(host, port))
	if err != nil {
		ready <- true
		done <- true
		t.Fatal(err)
	}

	ready <- true

	defer func() {
		ln.Close()
		done <- true
	}()

	for i := 0; i < count; i++ {
		conn, err := ln.Accept()
		if err != nil {
			t.Fatal(err)
		}
		conn.Close()
	}
}

func assertPingResult(t*testing.T, host, port string, expectSuccess bool, patterns ...string) {
	addr := net.JoinHostPort(host, port)
	err := Ping(testNetwork, addr, defaultTimeout)
	t.Logf("port ping %s -> %v", addr, err)

	if err != nil {
		if expectSuccess {
			t.Errorf("ping to %s failed; expected success", addr)
		} else {
			assertErrorContains(t, err, patterns...)
		}
	} else {
		if !expectSuccess {
			t.Errorf("ping to %s success; expected failure", addr)
		}
	}
}

func assertErrorContains(t*testing.T, err error, patterns ...string) {
	result := err.Error()
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
	if successCount != expectedSuccessCount {
		t.Errorf("expected %d successful pings, but got %d", expectedSuccessCount, successCount)
	}
}

func Test_ping_open_port(t*testing.T) {
	ready := make(chan bool)
	done := make(chan bool)
	go acceptN(t, testHost, testPort, 1, ready, done)
	<-ready

	assertPingResult(t, testHost, testPort, true)
	<-done

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
	done := make(chan bool)
	go acceptN(t, testHost, testPort, pingCount, ready, done)
	<-ready

	assertPingNSuccessCount(t, testHost, testPort, pingCount, pingCount)
	<-done
}

func Test_ping5_all_fail(t*testing.T) {
	pingCount := 5
	successCount := 0
	assertPingNSuccessCount(t, testHost, testPort, pingCount, successCount)
}

func Test_ping5_partial_success(t*testing.T) {
	successCount := 3
	ready := make(chan bool)
	done := make(chan bool)
	go acceptN(t, testHost, testPort, successCount, ready, done)
	<-ready

	pingCount := 5
	assertPingNSuccessCount(t, testHost, testPort, pingCount, successCount)
	<-done
}
