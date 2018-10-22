// Package main parses command line arguments and delegates to the correct functions.
package main

import (
	"controol/midi"
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
		return
	}

	var err error

	switch flag.Arg(0) {
	case "osc":
		err = oscCmd(flag.Args()[1:])
	case "midi":
		err = midiCmd(flag.Args()[1:])
	default:
		err = errors.New("unknown argument " + flag.Arg(0))
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func oscCmd(args []string) error {
	if len(args) == 0 {
		flag.Usage()
		return nil
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

func midiCmd(args []string) error {
	if len(args) == 0 {
		flag.Usage()
		return nil
	}

	switch args[0] {
	case "send":
		return midiSendCmd(args[1:])
	default:
		return errors.New("unknown argument " + args[0])
	}
	return nil
}

func midiSendCmd(args []string) error {
	if len(args) == 0 {
		flag.Usage()
		return nil
	}

	switch args[0] {
	case "cc":
		return midi.SendCC(args[1:])
	default:
		return errors.New("unknown argument " + args[0])
	}

	return nil
}
