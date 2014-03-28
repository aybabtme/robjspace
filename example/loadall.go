package main

import (
	"fmt"
	"github.com/aybabtme/rubyobj"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	trivialCommand = cli.Command{
		Name:        "trivial",
		Usage:       "decodes Ruby heap objects using the trivial decoder",
		Description: "The trivial decoder use a single core and reads a stream JSON objects with no particular delimiter.",
		Action:      loadTrivialAction("trivial"),
		Flags:       []cli.Flag{cli.StringFlag{Name: "filename", Usage: "name of the file to open"}},
	}

	parallelCommand = cli.Command{
		Name:        "parallel",
		Usage:       "decodes Ruby heap objects using a parallel decoder",
		Description: "The parallel decoder uses all cores. It expects newline-delimited JSON objects.",
		Action:      loadParallelAction("parallel"),
		Flags:       []cli.Flag{cli.StringFlag{Name: "filename", Usage: "name of the file to open"}},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "loadruby"
	app.Usage = "load a ruby heap dump file into memory"
	app.Commands = []cli.Command{trivialCommand, parallelCommand}

	app.Run(os.Args)
}

// Loader

func loadTrivialAction(command string) func(*cli.Context) {
	return func(c *cli.Context) {
		fr := loadFile(c, command)
		defer fr.Close()

		start := time.Now()

		dec := rubyobj.NewDecoder(fr)

		count := 0
		rObj := rubyobj.RubyObject{}

		var err error
		for {
			err = dec.Decode(&rObj)
			if err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				defer os.Exit(1)
				break
			}
			if err == io.EOF {
				break
			}
			count++
		}
		fmt.Printf("%d heap objects in %v\n", count, time.Since(start))
	}
}

func loadParallelAction(command string) func(*cli.Context) {
	return func(c *cli.Context) {
		defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(runtime.NumCPU()))

		fr := loadFile(c, command)
		defer fr.Close()

		start := time.Now()

		count := 0
		objC, errC := rubyobj.ParallelDecode(fr, uint(runtime.NumCPU()))

		wg := sync.WaitGroup{}
		wg.Add(1)
		go readErr(&wg, errC)

		for obj := range objC {
			_ = obj
			count++
		}

		wg.Wait()

		fmt.Printf("%d heap objects in %v\n", count, time.Since(start))
	}
}

func readErr(wg *sync.WaitGroup, errC <-chan error) {
	defer wg.Done()
	for err := range errC {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}

// Helpers

func loadFile(c *cli.Context, command string) *os.File {
	if !c.IsSet("filename") {
		cli.ShowCommandHelp(c, command)
		os.Exit(1)
	}

	fr, err := os.Open(c.String("filename"))
	fatal(err)

	fStat, err := fr.Stat()
	fatal(err)
	byteSize := fStat.Size()

	fmt.Printf("loading %s from '%s'\n", humanize.Bytes(uint64(byteSize)), fStat.Name())

	return fr
}

func fatal(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
