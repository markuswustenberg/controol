// Package osc has functions to send and receive OSC messages.
package osc

import (
	"strconv"

	osc2 "github.com/hypebeast/go-osc/osc"
	"github.com/pkg/errors"
)

// Send an OSC message to a specific host, port, and address.
// The argument slice should have host, port, and address as the first three entries,
// and values after that.
func Send(args []string) error {
	if len(args) < 4 {
		return errors.New("usage: osc send <host> <port> <address> <values>")
	}
	port, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.Wrap(err, "bad port "+args[1])
	}
	client := osc2.NewClient(args[0], port)
	msg := osc2.NewMessage(args[2])
	for _, arg := range args[3:] {
		if i, err := strconv.ParseInt(arg, 0, 32); err == nil {
			msg.Append(int32(i))
		} else if f, err := strconv.ParseFloat(arg, 32); err == nil {
			msg.Append(float32(f))
		} else if i, err := strconv.ParseInt(arg, 0, 64); err == nil {
			msg.Append(i)
		} else if f, err := strconv.ParseFloat(arg, 64); err == nil {
			msg.Append(f)
		} else if b, err := strconv.ParseBool(arg); err == nil {
			msg.Append(b)
		} else {
			msg.Append(arg)
		}
	}
	if err := client.Send(msg); err != nil {
		return errors.Wrap(err, "could not send message")
	}
	return nil
}

// Receive OSC messages at a specific host and port and print the message.
// Host should be a local interface.
// The callbacks are for testing currently.
func Receive(args []string, callbacks ...func(*osc2.Message)) error {
	if len(args) < 2 {
		return errors.New("usage: osc receive <host> <port>")
	}
	server := &osc2.Server{Addr: args[0] + ":" + args[1]}

	server.Handle("*", func(msg *osc2.Message) {
		osc2.PrintMessage(msg)
		for _, callback := range callbacks {
			callback(msg)
		}
	})

	if err := server.ListenAndServe(); err != nil {
		return errors.Wrap(err, "could not start receiving")
	}
	return nil
}
