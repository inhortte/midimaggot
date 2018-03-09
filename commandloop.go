package midimaggot

import (
	"bufio"
	"fmt"
	"os"
)

type doneThurks map[string](chan bool)

func makeDoneThurk(dts doneThurks, id string) doneThurks {
	dts[id] = make(chan bool, 1)
	return dts
}

func findDoneThurk(dts doneThurks, id string) chan bool {
	return dts[id]
}

func removeDoneThurk(dts doneThurks, id string) doneThurks {
	delete(dts, id)
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
	doneThurks := make(doneThurks, 128)
	sin := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("~> ")
		inp, _ := sin.ReadString('\n')
		for _, d := range commands {
			doneThurks = d.Thurk(inp, done, doneThurks)
		}
	}
}
