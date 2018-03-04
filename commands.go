package midimaggot

import (
	"fmt"
	"regexp"
	"strconv"
)

type directive interface {
	Thurk(string, chan<- bool, chan<- bool)
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

func (d *Done) Thurk(inp string, done chan<- bool, fRun chan<- bool) {
	parsed := parse(d.re, inp)
	if parsed != nil {
		done <- true
	}
	fRun <- true
}

type Bpm struct {
	re string
}

func (b *Bpm) Thurk(inp string, done chan<- bool, fRun chan<- bool) {
	parsed := parse(b.re, inp)
	if parsed != nil {
		bpm, err := strconv.Atoi(parsed[0])
		if err != nil {
			fmt.Println("Input problem, vole: ", err)
		} else {
			go sendMidiClock(bpm)
		}
	}
	fRun <- true
}
