package main

import (
	"net"
	"fmt"
	"regexp"
)

var pattern_getsockopt = regexp.MustCompile(`getsockopt: (.*)`)
var pattern_other = regexp.MustCompile(`^dial tcp: (.*)`)

func Ping(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)

	if err == nil {
		conn.Close()
	}
	return err
}

func PingN(host string, port int, count int, c chan error) {
	for i := 0; i < count; i++ {
		c <- Ping(host, port)
	}
}

func FormatResult(err error) string {
	if err == nil {
		return "success"
	}
	s := err.Error()
	if result := pattern_getsockopt.FindStringSubmatch(s); result != nil {
		return result[1]
	}
	if result := pattern_other.FindStringSubmatch(s); result != nil {
		return result[1]
	}
	return s
}

// TODO function to time the ping and return stats

// TODO functions to build total stats, aggregates
