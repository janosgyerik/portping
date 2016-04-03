portping
========

A simple library and command line utility to ping ports

[![GoDoc](https://godoc.org/github.com/janosgyerik/portping?status.svg)](https://godoc.org/github.com/janosgyerik/portping)
[![Build Status](https://travis-ci.org/janosgyerik/portping.svg?branch=master)](https://travis-ci.org/janosgyerik/portping)

Usage
-----

Ping port 80 of google.com 3 times:

    portping -c 3 google.com 80
    
Output:

    Starting to ping google.com:80 ...
    google.com:80 [1] -> success
    google.com:80 [2] -> success
    google.com:80 [3] -> success

Works with named ports too, for example:

    portping google.com http

See `portping -h` for all available options.

Download
--------

Binaries for several platforms are available on SourceForge:

https://sourceforge.net/projects/portping/files/

Generate test coverage report
-----------------------------

Run the commands:

    go test -coverprofile cover.out
    go tool cover -html=cover.out -o cover.html
    open cover.html

See more info: https://blog.golang.org/cover
