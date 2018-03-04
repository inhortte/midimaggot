package midimaggot

import (
	"bufio"
	"fmt"
	"os"
)

var cmdDone = Done{`^exit\s*$`}
var cmdBpm = Bpm{`^bpm\s+(\d+)\s*$`}
var commands = []directive{&cmdDone, &cmdBpm}

func CommandLoop(done chan bool) {
	sin := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("~> ")
		inp, _ := sin.ReadString('\n')
		for _, d := range commands {
			fRun := make(chan bool, 1)
			d.Thurk(inp, done, fRun)
		}
	}
}
