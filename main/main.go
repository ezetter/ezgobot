package main

import (
	"os"

	ezgobot "github.com/ezetter/ezgobot"
)

func main() {
	ezgobot.Init()
	ezgobot.ConversationLoop(os.Stdin, os.Stdout)
}
