package portping

import (
	"net"
	"fmt"
	"log"
)

func Ping(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)

	log.Printf("port ping %s -> %v", addr, err)

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

// TODO function to time the ping and return stats

// TODO functions to build total stats, aggregates
