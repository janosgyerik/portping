package main

// TODO
// args: host, port

// output:
// Port ping host:port
// Successful connection to host:port time=201ms
// Successful connection to host:port time=195ms
// Successful connection to host:port time=198ms
// --- host:port ping statistics ---
// n connections attempted, m successful, x% failed
// round-trip min/avg/max/stddev = a/b/c/d ms

// error output:
// portping: cannot resolve host: Unknown host
// portping: connect to address host: Connection refused

// flag: -c count
// flags: --tcp, --udp; default is tcp
// flag: -W timeout

// test with: 192.168.1.10
//tcp        0      0 127.0.1.1:53            0.0.0.0:*               LISTEN
//tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN
//tcp        0      0 127.0.0.1:631           0.0.0.0:*               LISTEN
//tcp        0      0 127.0.0.1:8089          0.0.0.0:*               LISTEN
//tcp6       0      0 :::22                   :::*                    LISTEN
//tcp6       0      0 ::1:631                 :::*                    LISTEN

// drop default count, print forever, until cancel with Control-C, and print stats

func main() {
}
