package main

import (
	"flag"
	"fmt"
	"github.com/aybabtme/rubyobj"
	"io"
	"os"
)

func main() {
	var filename string
	flag.StringVar(&filename, "file", "", "filename to read the ObjectSpace heap dump from")
	flag.Parse()

	var f *os.File
	if filename != "" {
		var err error
		f, err = os.Open(filename)
		if err != nil {
			perror(err)
		}
	} else {
		f = os.Stdin
	}

	rubyObj := rubyobj.RubyObject{}
	var err error

	for dec := rubyobj.NewDecoder(f); err == nil; err = dec.Decode(&rubyObj) {
		fmt.Printf("%v\n", &rubyObj)
	}

	if err != io.EOF {
		perror(err)
	}
}

func perror(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
