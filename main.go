// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	gflags "github.com/jessevdk/go-flags"
)

var version, commit, date = "unknown", "unknown", "unknown"

type Flags struct {
	Debug     bool   `short:"d" long:"debug" description:"enable debug output"`
	GeoDBPath string `long:"db" description:"path to read/store Maxmind DB"`
}

var flags Flags
var debug = log.New(ioutil.Discard, "", log.LstdFlags)

func main() {
	parser := gflags.NewParser(&flags, gflags.HelpFlag)
	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if flags.Debug {
		debug.SetOutput(os.Stdout)
	}

	catch()
	// Close things here.

	fmt.Println("exiting")
}

func catch() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	fmt.Println("listening for signal. CTRL+C to quit.")
	<-signals
	fmt.Println("\ninvoked termination, cleaning up")
}
