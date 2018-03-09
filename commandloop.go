package midimaggot

import (
	"bufio"
	"fmt"
	"os"
)

type doneThurk struct {
	id   string
	done chan bool
}

func makeDoneThurk(id string) doneThurk {
	dt := doneThurk{id, make(chan bool, 1)}
	return dt
}

func findDoneThurk(dts []doneThurk, id string) *doneThurk {
	for _, dt := range dts {
		if dt.id == id {
			return &dt
		}
	}
	return nil
}

func removeDoneThurk(dts *[]doneThurk, id string) *[]doneThurk {
	for idx, dt := range *dts {
		if dt.id == id {
			dts := append((*dts)[0:idx], (*dts)[idx+1:]...)
			return &dts
		}
	}
	return dts
}

var cmdDone = Done{`^exit\s*$`}
var cmdBpm = Bpm{`^bpm\s+(\d+)\s*$`}
var cmdProgramChange = ProgramChange{`^pc\s+(\d+)\s+(\d+)\s*$`}
var cmdPhaserIgnoreClock = PhaserIgnoreClock{`^pic\s+(\d+)\s*$`}
var cmdPhaserListenClock = PhaserListenClock{`^plc\s+(\d+)\s*$`}
var cmdPhaserBounceRate = PhaserBounceRate{`^pbr\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s*$`}
var cmdPhaserStopBounce = PhaserStopBounce{`^psb\s*$`}
var commands = []directive{&cmdDone, &cmdBpm, &cmdProgramChange, &cmdPhaserIgnoreClock,
	&cmdPhaserListenClock, &cmdPhaserBounceRate, &cmdPhaserStopBounce}

func CommandLoop(done chan bool) {
	doneThurks := make([]doneThurk, 0, 128)
	sin := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("~> ")
		inp, _ := sin.ReadString('\n')
		for _, d := range commands {
			dt := d.Thurk(inp, done, &doneThurks)
			doneThurks = append(doneThurks, dt)
		}
	}
}
