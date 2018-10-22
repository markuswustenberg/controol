// Package midi has functions to send and receive MIDI messages.
package midi

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/rakyll/portmidi"
)

func SendCC(args []string) error {
	if len(args) < 3 {
		return errors.New("usage: midi send cc <channel> <number> <value>")
	}

	channel, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("bad channel " + args[0])
	}
	if !(1 <= channel && channel <= 16) {
		return errors.Errorf("channel must be between 1 and 16, is %v", channel)
	}

	cc, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.New("bad cc number " + args[1])
	}
	if !(0 <= cc && cc <= 127) {
		return errors.Errorf("cc number must be between 0 and 127, is %v", cc)
	}

	v, err := strconv.Atoi(args[2])
	if err != nil {
		return errors.New("bad cc value " + args[2])
	}
	if !(0 <= v && v <= 127) {
		return errors.Errorf("cc value must be between 0 and 127, is %v", v)
	}

	if err := portmidi.Initialize(); err != nil {
		return errors.Wrap(err, "could not initialize portmidi")
	}
	defer portmidi.Terminate()

	id := portmidi.DefaultOutputDeviceID()
	out, err := portmidi.NewOutputStream(id, 1024, 0)
	if err != nil {
		return errors.Wrap(err, "could not create midi output")
	}
	defer out.Close()

	if err := out.WriteShort(int64(0xb0+channel-1), int64(cc), int64(v)); err != nil {
		return errors.Wrap(err, "could not send midi message")
	}

	return nil
}
