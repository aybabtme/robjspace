package main

import (
	"fmt"
	"github.com/aybabtme/rubyobj"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	"io"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "loadruby"
	app.Usage = "load a ruby heap dump file into memory"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "filename", Usage: "name of the file to open"},
	}
	app.Action = loadAction
	app.Run(os.Args)
}

func loadAction(c *cli.Context) {
	if !c.IsSet("filename") {
		cli.ShowAppHelp(c)
		return
	}

	fr, err := os.Open(c.String("filename"))
	fatal(err)
	defer fr.Close()

	fStat, err := fr.Stat()
	fatal(err)
	byteSize := fStat.Size()

	fmt.Printf("loading %s from '%s'\n", humanize.Bytes(uint64(byteSize)), fStat.Name())

	start := time.Now()

	dec := rubyobj.NewDecoder(fr)

	count := 0
	rObj := rubyobj.RubyObject{}

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

func fatal(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
