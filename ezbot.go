package ezbot

import (
	"bufio"
	"fmt"
	"io"
)

type state struct {
	transitions  map[string]state
	say          string
	memoryUpdate string
	memorySeek   []string
}

var currState, prevState state

var memory map[string]string

func buildState(say, memoryUpdate string, memoryRetrieve []string) state {
	newState := state{transitions: make(map[string]state),
		say: say, memoryUpdate: memoryUpdate, memorySeek: memoryRetrieve}
	return newState
}

// Init initializes the bot.
func Init() {
	memory = make(map[string]string)
	initialState := buildState("", "", nil)
	state2 := buildState("Hi. What's your name?", "name", nil)
	state3 := buildState("Hello %s! How can I help you?", "des_act", []string{"name"})
	state4 := buildState("Sorry %s, I don't know how to %s. How can I help you?", "des_act", []string{"name", "des_act"})
	initialState.transitions["default"] = state2
	state2.transitions["default"] = state3
	state3.transitions["default"] = state4
	state4.transitions["default"] = state4
	currState = state2
}

// ConversationLoop runs the bot's conversation loop.
func ConversationLoop(reader io.Reader, writer io.Writer) {
	newReader := bufio.NewReader(reader)
	input := ""
	for {
		fmt.Fprint(writer, act(input))
		input, _ = newReader.ReadString('\n')
		if input == "exit" {
			return
		}
	}
}

func act(input string) string {
	if input != "" {
		memory[prevState.memoryUpdate] = input[:len(input)-1]
	}
	memoryOut := retrieveMemory(currState)
	out := fmt.Sprintln()
	out += fmt.Sprintf(currState.say, memoryOut...)

	out += fmt.Sprint(" ")
	prevState = currState
	currState = currState.transitions["default"]
	return out
}

func retrieveMemory(fromState state) []interface{} {
	memoryOut := make([]interface{}, len(currState.memorySeek))
	for i, str := range currState.memorySeek {
		memoryOut[i] = memory[str]
	}
	return memoryOut
}
