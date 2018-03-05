package midimaggot

import (
	"fmt"
	"github.com/rakyll/portmidi"
	"time"
)

type Control byte

const empressPhaserChannel int64 = 7
const brothersChannel int64 = 8
const gravitasChannel int64 = 8
const (
	TimingClock   Control = 0xf8
	ControlChange Control = 0xb0
)

func devices() (map[string]portmidi.DeviceID, map[string]portmidi.DeviceID) {
	numDevices := portmidi.CountDevices()
	fmt.Println("MIDI Devices: ", numDevices)
	inputs := make(map[string]portmidi.DeviceID, numDevices)
	outputs := make(map[string]portmidi.DeviceID, numDevices)
	for i := 0; i < numDevices; i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		if info.IsInputAvailable {
			inputs[info.Name] = portmidi.DeviceID(i)
		}
		if info.IsOutputAvailable {
			outputs[info.Name] = portmidi.DeviceID(i)
		}
	}
	fmt.Println("Inputs: ", inputs)
	fmt.Println("Outputs: ", outputs)
	return inputs, outputs
}
func getpisoundInStream() *portmidi.Stream {
	inputs, _ := devices()
	inputId := inputs["pisound MIDI PS-225TT43"]
	fmt.Println("inputID: ", inputId)
	inp, err := portmidi.NewInputStream(inputId, 1024)
	if err != nil {
		fmt.Println("Cannot initialize input stream: ", err)
	}
	return inp
}
func getpisoundOutStream() *portmidi.Stream {
	_, outputs := devices()
	outputId := outputs["pisound MIDI PS-225TT43"]
	fmt.Println("outputID: ", outputId)
	out, err := portmidi.NewOutputStream(outputId, 1024, 0)
	if err != nil {
		fmt.Println("Cannot initialize output stream: ", err)
	}
	return out
}

/*
 * Midi divides a quarter into 24 ticks that the clock individually sends.
 * Thus, the arithmetic.
 */
func sendMidiClock(bpm int) {
	sleepTime := (60.0 / float64(bpm)) / 24.0
	fmt.Printf("bpm: %v, sleep time: %v\n", bpm, sleepTime)
	out := getpisoundOutStream()
	for i := 0; i < 24; i++ {
		out.WriteShort(0xf8, 0, 100)
		time.Sleep(time.Duration(sleepTime * 1000000000))
	}
	out.Close()
}

func sendProgramChange(channel, program int) {
	out := getpisoundOutStream()
	var eType int64 = 0xc0 | int64(channel-1)
	out.WriteShort(eType, int64(program), 0)
	out.Close()
}

func empressPhaserIgnoreClock(args ...int) {
	channel := empressPhaserChannel
	if len(args) > 0 {
		channel = int64(args[0])
	}
	out := getpisoundOutStream()
	controlChange := int64(0xb0 & (channel - 1))
	controlNumber := int64(51)
	controlValue := int64(0)
	out.WriteShort(controlChange, controlNumber, controlValue)
	out.Close()
}

func empressPhaserListenClock(args ...int) {
	channel := empressPhaserChannel
	if len(args) > 0 {
		channel = int64(args[0])
	}
	out := getpisoundOutStream()
	defer out.Close()
	controlChange := int64(0xb0 & (channel - 1))
	controlNumber := int64(51)
	controlValue := int64(127)
	out.WriteShort(controlChange, controlNumber, controlValue)
}

func ProgramChangeForward() {
	fmt.Println("starting ProgramChangeForward")
	in := getpisoundInStream()
	defer in.Close()
	ch := in.Listen()
	for {
		event := <-ch
		eType := event.Status >> 4
		channel := event.Status % 16
		program := event.Data1
		if eType == 12 {
			fmt.Printf("Forwarding Program Change to program %v from channel %v to channel 8, because the Brothers and Gravitas are listening on EIGHT (that's 7 for me), vole!", program, channel)
		}
		out := getpisoundOutStream()
		out.WriteShort(0xc7, program, 0)
		out.Close()
	}
}
