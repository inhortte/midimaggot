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

var cmdUsage = Usage{`^(usage|help)\s*$`}
var cmdDone = Done{`^exit\s*$`}
var cmdBpm = Bpm{`^bpm\s+(\d+)\s*$`}
var cmdProgramChange = ProgramChange{`^pc\s+(\d+)\s+(\d+)\s*$`}
var cmdPhaserIgnoreClock = PhaserIgnoreClock{`^pic\s+(\d+)\s*$`}
var cmdPhaserListenClock = PhaserListenClock{`^plc\s+(\d+)\s*$`}
var cmdPhaserBounceRate = PhaserBounceRate{`^pbr\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s*$`}
var cmdPhaserStopBounce = PhaserStopBounce{`^psb\s*$`}
var cmdPhaserRandomRate = PhaserRandomRate{`^prr\s+(\d+)\s+(\d+)\s+(\d+)\s*$`}
var cmdPhaserStopRandomRate = PhaserStopRandomRate{`^psrr\s*$`}
var cmdProgram = Program{`^p\s+(\d+)\s+"(.+)"\s*$`}
var commands = []directive{&cmdUsage, &cmdDone, &cmdBpm, &cmdProgramChange, &cmdPhaserIgnoreClock,
	&cmdPhaserListenClock, &cmdPhaserBounceRate, &cmdPhaserStopBounce,
	&cmdPhaserRandomRate, &cmdPhaserStopRandomRate, &cmdProgram}

func processCommand(inp string, done chan<- bool, dts doneThurks) doneThurks {
	for _, d := range commands {
		dts = d.Thurk(inp, done, dts)
	}
	return dts
}

func CommandLoop(done chan bool) {
	doneThurks := make(doneThurks, 128)
	go ProgramChangeForward(done, doneThurks)
	sin := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("~> ")
		inp, _ := sin.ReadString('\n')
		doneThurks = processCommand(inp, done, doneThurks)
		/*
			for _, d := range commands {
				doneThurks = d.Thurk(inp, done, doneThurks)
			}
		*/
	}
}
