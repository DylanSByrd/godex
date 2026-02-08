package main

import (
	"strings"
	"fmt"
	"os"
	"bufio"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func startRepl(context *commandContext) {
	scanner := bufio.NewScanner(os.Stdin)
	cmds := getCommands()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		input := scanner.Text()
		tokens := cleanInput(input)

		if len(tokens) == 0 {
			continue
		}

		cmd, ok := cmds[tokens[0]]
		args := tokens[1:]
		if ok {
			err := cmd.callback(context, args...)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Uknown command")
		}
	}
}
