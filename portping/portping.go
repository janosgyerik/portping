package portping

import (
	"net"
	"fmt"
)

func Ping(host string, port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err == nil {
		conn.Close()
	}
	return err
}

// TODO function to ping repeatedly

// TODO function to time the ping and return stats

// TODO functions to build total stats, aggregates
