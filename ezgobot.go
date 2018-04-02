package ezgobot

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type transitionMapping struct {
	inputMatch *regexp.Regexp
	transition string
}

// State defines a node in the state machine representing a conversation.
type State struct {
	id                      string
	transitions             map[string]*State
	immediateTransition     *State
	transitionInputMappings []transitionMapping
	say                     string
	memoryWrite             string
	memoryRead              []string
}

var currState State

var memory map[string]string

var debug bool

// BuildTransitionState builds a state as a transition from the reciever
func (s *State) BuildTransitionState(say string, transitionName string) *State {
	newState := &State{transitions: make(map[string]*State), say: say}
	s.transitions[transitionName] = newState
	return newState
}

// AddTransitionMapping adds a mapping from inputs to transitions.
func (s *State) AddTransitionMapping(matchRegExp, transition string) *State {
	inputMatch, _ := regexp.Compile(matchRegExp)
	s.transitionInputMappings = append(s.transitionInputMappings, transitionMapping{inputMatch, transition})
	return s
}

// AddImmediateTransition adds a transition that is immediate, without
// requiring inputs. This allows reuse of behaviors from other states.
func (s *State) AddImmediateTransition(transitionState *State) *State {
	s.immediateTransition = transitionState
	return s
}

// SetID sets the state's identifier.
func (s *State) SetID(id string) *State {
	s.id = id
	return s
}

// SetMemoryWrite sets the memory that will be written after the state runs.
func (s *State) SetMemoryWrite(memoryWrite string) *State {
	s.memoryWrite = memoryWrite
	return s
}

// SetMemoryRead sets the memory elements that the state needs to read.
func (s *State) SetMemoryRead(memoryRead []string) *State {
	s.memoryRead = memoryRead
	return s
}

// AddTransition adds a new state transition.
func (s *State) AddTransition(transitionName string, destination *State) *State {
	s.transitions[transitionName] = destination
	return s
}

// Init initializes the bot.
func Init() {
	memory = make(map[string]string)
	currState = State{id: "boot", transitions: make(map[string]*State)}
	s1 := currState.BuildTransitionState("Hi. What's your name?", "default").
		SetMemoryWrite("name").SetID("init")
	s2 := s1.BuildTransitionState("Hello %s! How can I help you?", "default").
		SetID("helpful").
		SetMemoryWrite("des_act").
		SetMemoryRead([]string{"name"})
	s2.BuildTransitionState("Sorry %s, I don't know how to %s. How can I help you?", "default").
		SetID("cantdo").
		SetMemoryRead([]string{"name", "des_act"}).
		AddImmediateTransition(s2)

	s2.BuildTransitionState("My name is Machina. How can I help you?", "ask_name").
		SetID("ask_name").
		AddImmediateTransition(s2)
	s2.AddTransitionMapping("what is your name", "ask_name")
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

func normalizeInput(input string) string {
	re := regexp.MustCompile("[\\?\\.!]")
	return re.ReplaceAllString(strings.ToLower(input), "")
}

func determineTransition(input string, mappings []transitionMapping) string {
	for _, mapping := range mappings {
		if mapping.inputMatch.MatchString(input) {
			return mapping.transition
		}
	}
	return "default"
}

func printDebugIn(input string) {
	fmt.Printf("input: %s\n", input)
	fmt.Printf("  Curr state: %s\n", currState.id)
}

func printDebugOut(output string) {
	fmt.Printf("  Curr state: %s\n", currState.id)

	fmt.Printf("  Output: %s\n", output)
}
func act(input string) string {
	if debug {
		printDebugIn(input)
	}
	if input != "" {
		memory[currState.memoryWrite] = input[:len(input)-1]
	}
	nextTransition := determineTransition(normalizeInput(input), currState.transitionInputMappings)

	out := fmt.Sprintln()
	if _, ok := currState.transitions[nextTransition]; ok {
		currState = *currState.transitions[nextTransition]
		remembered := retrieveMemory(currState)
		out += fmt.Sprintf(currState.say, remembered...)
	} else {
		// TODO: this needs to be impossible.
		out += "I'm confused! I don't know what to do!"
	}
	out += fmt.Sprint(" ")
	if currState.immediateTransition != nil {
		currState = *currState.immediateTransition
	}
	if debug {
		printDebugOut(out)
	}
	return out
}

func retrieveMemory(fromState State) []interface{} {
	memoryOut := make([]interface{}, len(currState.memoryRead))
	for i, str := range currState.memoryRead {
		memoryOut[i] = memory[str]
	}
	return memoryOut
}
