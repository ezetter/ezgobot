package main

import (
	"os"

	"github.com/ezetter/ezbot"
)

func main() {
	ezbot.Init()
	ezbot.ConversationLoop(os.Stdin, os.Stdout)
}
