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

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	cmds := getCommands()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		input := scanner.Text()
		tokens := cleanInput(input)

		if len(tokens) != 0 {
			cmd, ok := cmds[tokens[0]]
			if ok {
				err := cmd.callback()
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Uknown command")
			}
		}
	}
}
