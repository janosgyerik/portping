package portping

import (
	"net"
	"fmt"
)

// TODO return the error
func Ping(host string, port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	fmt.Printf("%s:%d %v", host, port, err)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

// TODO function to ping repeatedly

// TODO function to time the ping and return stats

// TODO functions to build total stats, aggregates
