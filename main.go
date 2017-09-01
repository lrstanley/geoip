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
	"sync"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	gflags "github.com/jessevdk/go-flags"
)

var version, commit, date = "unknown", "unknown", "unknown"

type Bot struct{}

type Flags struct {
	ConfigFile string `short:"c" long:"config" description:"configuration file location" default:"config.toml"`
	Debug      bool   `short:"d" long:"debug" description:"enable debug output"`
	GenConfig  bool   `long:"gen-config" description:"generate and output an example configuration file"`
}

type Config struct {
	DatabasePath   string
	ReconnectDelay int               `toml:"reconnect_delay"`
	Servers        map[string]Server `toml:"servers"`
}

var flags Flags
var conf *Config
var debug = log.New(ioutil.Discard, "", log.LstdFlags)

func main() {
	parser := gflags.NewParser(&flags, gflags.HelpFlag)
	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if flags.GenConfig {
		genConfig()
		return
	}

	if flags.Debug {
		debug.SetOutput(os.Stdout)
	}

	conf = &Config{
		DatabasePath:   "state.db",
		ReconnectDelay: 45,
	}

	_, err = toml.DecodeFile(flags.ConfigFile, conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config[%s]: %s\n", flags.ConfigFile, err)
		os.Exit(1)
	}

	if conf.ReconnectDelay < 10 {
		conf.ReconnectDelay = 10
	}

	if len(conf.Servers) == 0 {
		fmt.Fprintf(os.Stderr, "config[%s]: no servers specified", flags.ConfigFile)
		os.Exit(1)
	}

	done := make(chan struct{})
	var wg sync.WaitGroup

	for key, _ := range conf.Servers {
		time.Sleep(1 * time.Second)

		go func(server Server, id string) {
			wg.Add(1)
			server.setup(id, done)
			wg.Done()
		}(conf.Servers[key], key)
	}

	catch()
	close(done)
	wg.Wait()

	fmt.Println("exiting")
}

func catch() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	fmt.Println("listening for signal. CTRL+C to quit.")
	<-signals
	fmt.Println("\ninvoked termination, cleaning up")
}

func genConfig() {
	fmt.Fprint(os.Stdout, `# example\foo = "bar"\n\n`)
}
