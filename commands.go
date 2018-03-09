package midimaggot

import (
	"fmt"
	"regexp"
	"strconv"
)

type directive interface {
	Thurk(string, chan<- bool, *[]doneThurk) doneThurk
}

func parse(reString string, inp string) []string {
	re, err := regexp.Compile(reString)
	if err != nil {
		panic(err)
	}
	reMatch := re.FindStringSubmatch(inp)
	if reMatch != nil {
		fmt.Println("parsed: ", reMatch)
		return reMatch[1:]
	} else {
		return nil
	}
}

type Done struct {
	re string
}

func (d *Done) Thurk(inp string, done chan<- bool, dts *[]doneThurk) doneThurk {
	parsed := parse(d.re, inp)
	if parsed != nil {
		done <- true
	}
	return makeDoneThurk("done")
}

type Bpm struct {
	re string
}

func (b *Bpm) Thurk(inp string, done chan<- bool, dts *[]doneThurk) doneThurk {
	parsed := parse(b.re, inp)
	if parsed != nil {
		bpm, err := strconv.Atoi(parsed[0])
		if err != nil {
			fmt.Println("Input problem, vole: ", err)
		} else {
			go sendMidiClock(bpm)
		}
	}
	return makeDoneThurk("bpm")
}

type ProgramChange struct {
	re string
}

func (pc *ProgramChange) Thurk(inp string, done chan<- bool, dts *[]doneThurk) doneThurk {
	parsed := parse(pc.re, inp)
	if parsed != nil {
		channel, err := strconv.Atoi(parsed[0])
		program, err2 := strconv.Atoi(parsed[1])
		if err != nil || err2 != nil {
			fmt.Println("Input problem, vole: ", err)
		} else {
			go sendProgramChange(channel, program)
		}
	}
	return makeDoneThurk("program-change")
}

type PhaserIgnoreClock struct {
	re string
}

func (pic *PhaserIgnoreClock) Thurk(inp string, done chan<- bool, dts *[]doneThurk) doneThurk {
	parsed := parse(pic.re, inp)
	if parsed != nil {
		channel, err := strconv.Atoi(parsed[0])
		if err != nil {
			fmt.Println("Input problem, honey: ", err)
		} else {
			go empressPhaserIgnoreClock(channel)
		}
	}
	return makeDoneThurk("phaser-ignore-clock")
}

type PhaserListenClock struct {
	re string
}

func (plc *PhaserListenClock) Thurk(inp string, done chan<- bool, dts *[]doneThurk) doneThurk {
	parsed := parse(plc.re, inp)
	if parsed != nil {
		channel, err := strconv.Atoi(parsed[0])
		if err != nil {
			fmt.Println("Input problem, honey: ", err)
		} else {
			go empressPhaserListenClock(channel)
		}
	}
	return makeDoneThurk("phaser-listen-clock")
}

type PhaserBounceRate struct {
	re string
}

func (pbr *PhaserBounceRate) Thurk(inp string, done chan<- bool, dts *[]doneThurk) doneThurk {
	dt := makeDoneThurk("phaser-bounce-rate")
	parsed := parse(pbr.re, inp)
	if parsed != nil {
		channel, err1 := strconv.Atoi(parsed[0])
		bpm, err2 := strconv.Atoi(parsed[1])
		low, err3 := strconv.Atoi(parsed[2])
		high, err4 := strconv.Atoi(parsed[3])
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			fmt.Println("Input problem: pbr channel bpm low high ... ")
		} else {
			go empressPhaserBounceRate(&dt, channel, bpm, low, high)
		}
	}
	return dt
}

type PhaserStopBounce struct {
	re string
}

func (psb *PhaserStopBounce) Thurk(inp string, done chan<- bool, dts *[]doneThurk) doneThurk {
	parsed := parse(psb.re, inp)
	if parsed != nil {
		dt := findDoneThurk(*dts, "phaser-bounce-rate")
		dt.done <- true
		dts = removeDoneThurk(dts, "phaser-bounce-rate")
	}
	return makeDoneThurk("phaser-stop-bounnce")
}
