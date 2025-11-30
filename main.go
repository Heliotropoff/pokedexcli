package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Print("Pokedex > ")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		new_input := scanner.Text()
		parsed_input := cleanInput(new_input)
		//fmt.Printf("Your command was: %s\n", parsed_input[0])
		if len(parsed_input) > 0 {
			for _, cmdName := range parsed_input {
				cmd, ok := supportedCommands[cmdName]
				if !ok {
					fmt.Println("Unknown command")
				} else {
					err := cmd.callback(&config)
					if err != nil {
						fmt.Errorf("Error occurred: %w", err)
					}
				}
			}
			fmt.Print("Pokedex > ")
		}
	}
}
