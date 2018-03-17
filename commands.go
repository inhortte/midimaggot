package midimaggot

import (
	"fmt"
	"regexp"
	"strconv"
)

type directive interface {
	Thurk(string, chan<- bool, doneThurks) doneThurks
}

func parse(reString string, inp string) []string {
	re, err := regexp.Compile(reString)
	if err != nil {
		panic(err)
	}
	reMatch := re.FindStringSubmatch(inp)
	if reMatch != nil {
		// fmt.Println("parsed: ", reMatch)
		return reMatch[1:]
	} else {
		return nil
	}
}

type Done struct {
	re string
}

func (d *Done) Thurk(inp string, done chan<- bool, dts doneThurks) doneThurks {
	parsed := parse(d.re, inp)
	if parsed != nil {
		done <- true
	}
	return dts
}

type Usage struct {
	re string
}

func (u *Usage) Thurk(inp string, done chan<- bool, dts doneThurks) doneThurks {
	usage := map[string]string{
		"usage": "Display a summary of the various maggot types: usage | help",
		"exit":  "Extract yourself from the maggot: exit",
		"bpm":   "Send midi clock: bpm <beats-per-minute>",
		"pc":    "Program change: pc <channel> <program>",
		"pic":   "Empress phaser ignore clock: pic",
		"plc":   "Empress phaser listen clock: plc",
		"pbr":   "Empress phaser bounce rate: pbr <channel> <beats-per-minute> <low> <high>",
		"psd":   "Empress phaser stop bounce rate: psb",
		"prr":   "Empress phaser random rate: prr <channel> <beats-per-minute> <division>",
		"psrr":  "Empress phaser stop random rate: psrr",
	}
	parsed := parse(u.re, inp)
	if parsed != nil {
		fmt.Println()
		for k, v := range usage {
			fmt.Printf("%s -- %s\n", k, v)
		}
	}
	return dts
}

type Bpm struct {
	re string
}

func (b *Bpm) Thurk(inp string, done chan<- bool, dts doneThurks) doneThurks {
	parsed := parse(b.re, inp)
	if parsed != nil {
		bpm, err := strconv.Atoi(parsed[0])
		if err != nil {
			fmt.Println("Input problem, vole: ", err)
		} else {
			go sendMidiClock(bpm)
		}
	}
	return dts
}

type ProgramChange struct {
	re string
}

func (pc *ProgramChange) Thurk(inp string, done chan<- bool, dts doneThurks) doneThurks {
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
	return dts
}

type PhaserIgnoreClock struct {
	re string
}

func (pic *PhaserIgnoreClock) Thurk(inp string, done chan<- bool, dts doneThurks) doneThurks {
	parsed := parse(pic.re, inp)
	if parsed != nil {
		channel, err := strconv.Atoi(parsed[0])
		if err != nil {
			fmt.Println("Input problem, honey: ", err)
		} else {
			go empressPhaserIgnoreClock(channel)
		}
	}
	return dts
}

type PhaserListenClock struct {
	re string
}

func (plc *PhaserListenClock) Thurk(inp string, done chan<- bool, dts doneThurks) doneThurks {
	parsed := parse(plc.re, inp)
	if parsed != nil {
		channel, err := strconv.Atoi(parsed[0])
		if err != nil {
			fmt.Println("Input problem, honey: ", err)
		} else {
			go empressPhaserListenClock(channel)
		}
	}
	return dts
}

type PhaserBounceRate struct {
	re string
}

func (pbr *PhaserBounceRate) Thurk(inp string, done chan<- bool, dts doneThurks) doneThurks {
	// dts = makeDoneThurk(dts, "phaser-bounce-rate")
	parsed := parse(pbr.re, inp)
	if parsed != nil {
		channel, err1 := strconv.Atoi(parsed[0])
		bpm, err2 := strconv.Atoi(parsed[1])
		low, err3 := strconv.Atoi(parsed[2])
		high, err4 := strconv.Atoi(parsed[3])
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			fmt.Println("Input problem: pbr channel bpm low high ... ")
		} else {
			dts["phaser-bounce-rate"] = make(chan bool, 1)
			go empressPhaserBounceRate(dts["phaser-bounce-rate"], channel, bpm, low, high)
		}
	}
	return dts
}

type PhaserStopBounce struct {
	re string
}

func (psb *PhaserStopBounce) Thurk(inp string, done chan<- bool, dts doneThurks) doneThurks {
	parsed := parse(psb.re, inp)
	if parsed != nil {
		fmt.Println("phaser stop bounce....")
		// c := findDoneThurk(dts, "phaser-bounce-rate")
		//c <- true
		dts["phaser-bounce-rate"] <- true
		// dts = removeDoneThurk(dts, "phaser-bounce-rate")
		delete(dts, "phaser-bounce-rate")
	}
	return dts
}

type PhaserRandomRate struct {
	re string
}

func (prr *PhaserRandomRate) Thurk(inp string, done chan<- bool, dts doneThurks) doneThurks {
	parsed := parse(prr.re, inp)
	if parsed != nil {
		channel, err1 := strconv.Atoi(parsed[0])
		bpm, err2 := strconv.Atoi(parsed[1])
		division, err3 := strconv.ParseFloat(parsed[2], 64)
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Println("Input problem prr channel bpm division ... ")
		} else {
			dts["phaser-random-rate"] = make(chan bool, 1)
			go empressPhaserRandomRate(dts["phaser-random-rate"], channel, bpm, division)
		}
	}
	return dts
}

type PhaserStopRandomRate struct {
	re string
}

func (psrr *PhaserStopRandomRate) Thurk(inp string, done chan<- bool, dts doneThurks) doneThurks {
	parsed := parse(psrr.re, inp)
	if parsed != nil {
		fmt.Println("phaser stop random rate ...")
		dts["phaser-random-rate"] <- true
		delete(dts, "phaser-random-rate")
	}
	return dts
}
