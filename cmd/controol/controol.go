// Package main parses command line arguments and delegates to the correct functions.
package main

import (
	"controol/osc"
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}

	var err error

	switch flag.Arg(0) {
	case "osc":
		err = oscCmd(flag.Args()[1:])
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func oscCmd(args []string) error {
	if len(args) == 0 {
		flag.Usage()
	}

	switch args[0] {
	case "send":
		return osc.Send(args[1:])
	case "receive":
		return osc.Receive(args[1:])
	default:
		return errors.New("unknown argument " + args[0])
	}
	return nil
}
