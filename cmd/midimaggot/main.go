package main

import (
	"github.com/inhortte/midimaggot"
	"github.com/rakyll/portmidi"
)

func main() {
	mainDone := make(chan bool, 1)

	portmidi.Initialize()
	defer portmidi.Terminate()
	go midimaggot.ProgramChangeForward()
	go midimaggot.CommandLoop(mainDone)
	<-mainDone
}
