# controol

[![Build Status](https://travis-ci.org/wustenberg/controol.svg?branch=master)](https://travis-ci.org/wustenberg/controol)

`controol` is a small control tool to work with OSC and MIDI messages.

## Usage

In one terminal, try

    controol osc receive localhost 9000

and then in another

    controol osc send localhost 9000 /tentoonstelling 1 2.0 true hurray
