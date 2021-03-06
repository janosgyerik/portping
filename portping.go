/*
Package portping provides simple functions to ping TCP ports.
It also includes a simple command line interface.
 */
package portping

import (
	"net"
	"time"
)

// Ping connects to the address on the named network,
// using net.DialTimeout, and immediately closes it.
// It returns the connection error. A nil value means success.
// For examples of valid values of network and address,
// see the documentation of net.Dial
func Ping(network, address string, timeout time.Duration) error {
	conn, err := net.DialTimeout(network, address, timeout)
	if conn != nil {
		defer conn.Close()
	}
	return err
}

// PingN calls Ping the specified number of times,
// and sends the results to the given channel.
func PingN(network, address string, timeout time.Duration, count int, c chan <- error) {
	for i := 0; i < count; i++ {
		c <- Ping(network, address, timeout)
	}
}
