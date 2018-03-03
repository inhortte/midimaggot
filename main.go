package main

import (
	"bufio"
	"fmt"
	"github.com/rakyll/portmidi"
	"os"
	"regexp"
	"strconv"
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

func programChangeForward() {
	in := getpisoundInStream()
	defer in.Close()
	ch := in.Listen()
	for {
		event := <-ch
		fmt.Println("received -> ", event)
	}
}

func commandLoop(done chan bool) {
	sin := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("~> ")
		inp, _ := sin.ReadString('\n')
		cmdDone, _ := regexp.Compile(`^exit\s*$`)
		cmdBpm, _ := regexp.Compile(`^bpm\s+(\d+)\s*$`)
		if cmdDone.Match([]byte(inp)) {
			done <- true
		}
		bpmMatch := cmdBpm.FindStringSubmatch(inp)
		if bpmMatch != nil {
			bpm, _ := strconv.Atoi(bpmMatch[1])
			go sendMidiClock(bpm)
		}
	}
}

func main() {
	mainDone := make(chan bool, 1)

	portmidi.Initialize()
	defer portmidi.Terminate()
	go programChangeForward()
	// sendMidiClock(90)
	go commandLoop(mainDone)
	<-mainDone
}
