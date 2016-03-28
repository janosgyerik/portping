package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// TODO
// flags: --tcp, --udp; default is tcp
// flag: -W timeout
// drop default count, print forever, until cancel with Control-C, and print stats

const defaultCount = 5

func exit() {
	flag.Usage()
	os.Exit(1)
}

type Params struct {
	host  string
	port  int
	count int
}

func parseArgs() Params {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] host port\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	countPtr := flag.Int("count", defaultCount, "stop after count connections")
	flag.Parse()

	if len(flag.Args()) < 2 {
		exit()
	}

	host := flag.Args()[0]
	port, parseErr := strconv.Atoi(flag.Args()[1])
	if parseErr != nil {
		exit()
	}

	return Params{
		host: host,
		port: port,
		count: *countPtr,
	}
}

func main() {
	params := parseArgs()

	host := params.host
	port := params.port
	count := params.count

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Starting to ping %s ...\n", addr)

	c := make(chan error)
	go PingN(host, port, count, c)

	allSuccessful := true

	for i := 0; i < count; i++ {
		// TODO print details only if verbose, otherwise print just OpError.Err
		var msg string
		if err := <-c; err == nil {
			msg = "success"
		} else {
			msg = err.Error()
			allSuccessful = false
		}
		// TODO add time
		fmt.Printf("port ping %s [%d] -> %s\n", addr, i + 1, msg)
	}

	// TODO print summary
	// --- host:port ping statistics ---
	// n connections attempted, m successful, x% failed
	// round-trip min/avg/max/stddev = a/b/c/d ms

	if !allSuccessful {
		os.Exit(1)
	}
}
