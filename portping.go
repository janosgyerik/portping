/*
Package portping provides simple functions to ping TCP ports.
It also includes a simple command line interface.
 */
package portping

import (
	"net"
	"time"
)

const defaultTimeout = 10 * time.Second

// Ping connects to the specified host and port
// using net.DialTimeout and network "tcp".
func Ping(host, port string) error {
	addr := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", addr, defaultTimeout)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		return err
	}
	return nil
}

// PingN calls Ping the specified number of times,
// and sends the results to the given channel.
func PingN(host, port string, count int, c chan error) {
	for i := 0; i < count; i++ {
		c <- Ping(host, port)
	}
}

// FormatResult converts the result returned by Ping to string.
func FormatResult(err error) string {
	if err == nil {
		return "success"
	}
	switch err := err.(type) {
	case *net.OpError:
		return err.Err.Error()
	default:
		return err.Error()
	}
}
