// Command line interface to ping ports
package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/janosgyerik/portping"
	"net"
	"time"
)

// TODO
// flags: --tcp, --udp; default is tcp
// flag: -v verbose; default=false
// drop default count, print forever, until cancel with Control-C, and print stats

const (
	defaultCount = 5
	defaultTimeout = 10
	defaultNetwork = "tcp"
)

func exit() {
	flag.Usage()
	os.Exit(1)
}

type Params struct {
	host    string
	port    string
	count   int
	timeout time.Duration
}

func parseArgs() Params {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] host port\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	countPtr := flag.Int("c", defaultCount, "stop after count connections")
	timeoutPtr := flag.Int("W", defaultTimeout, "time in seconds to wait for connections")
	flag.Parse()

	if len(flag.Args()) < 2 {
		exit()
	}

	host := flag.Args()[0]
	port := flag.Args()[1]

	return Params{
		host: host,
		port: port,
		count: *countPtr,
		timeout: time.Duration(*timeoutPtr) * time.Second,
	}
}


// FormatResult converts the result returned by Ping to string.
func formatResult(err error) string {
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

func main() {
	params := parseArgs()

	host := params.host
	port := params.port
	count := params.count

	addr := net.JoinHostPort(host, port)
	fmt.Printf("Starting to ping %s ...\n", addr)

	c := make(chan error)
	go portping.PingN(defaultNetwork, addr, params.timeout, count, c)

	allSuccessful := true

	for i := 0; i < count; i++ {
		// TODO add time
		err := <-c
		if err != nil {
			allSuccessful = false
		}
		fmt.Printf("%s [%d] -> %s\n", addr, i + 1, formatResult(err))
	}

	// TODO print summary
	// --- host:port ping statistics ---
	// n connections attempted, m successful, x% failed
	// round-trip min/avg/max/stddev = a/b/c/d ms

	if !allSuccessful {
		os.Exit(1)
	}
}
