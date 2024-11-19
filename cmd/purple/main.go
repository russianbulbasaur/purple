package main

import (
	"fmt"
	"os"
	"purple/internals/server"
	"strings"
)

func main() {
	args := os.Args
	var argsMap map[string]string
	if len(args) > 1 {
		argsMap = parseArguments(args)
	}
	purpleServer := server.NewServer(8000, "localhost", argsMap)
	purpleServer.Listen()
}

func parseArguments(args []string) map[string]string {
	argsMap := make(map[string]string)
	for index := 1; index < len(args); index += 2 {
		argument := strings.ToLower(args[index])
		if isValidArg(argument) {
			argument = strings.TrimPrefix(argument, "--")
			switch argument {
			case "dir":
				if index+1 < len(args) && isValidValue(args[index+1]) {
					argsMap[argument] = args[index+1]
					continue
				}
			case "dbfilename":
				if index+1 < len(args) && isValidValue(args[index+1]) {
					argsMap[argument] = args[index+1]
					continue
				}
			}
		}
		panic(fmt.Sprintf("Invalid argument %s", argument))
	}
	return argsMap
}

func isValidArg(arg string) bool {
	return strings.HasPrefix(arg, "--")
}

func isValidValue(value string) bool {
	return !strings.HasPrefix(value, "--") && (strings.TrimSpace(value) != "")
}
