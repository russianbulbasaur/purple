package main

import (
	"fmt"
	"log"
	"os"
	"purple/internals/server"
	"strconv"
	"strings"
)

func main() {
	args := os.Args
	var argsMap map[string]interface{} = make(map[string]interface{})
	if len(args) > 1 {
		parseArguments(args, argsMap)
	}
	checkArguments(argsMap)
	port, err := strconv.Atoi(argsMap["port"].(string))
	if err != nil {
		log.Fatalln("Invalid port")
	}
	purpleServer := server.NewServer(port, "localhost", argsMap)
	purpleServer.Listen()
}

func parseArguments(args []string, argsMap map[string]interface{}) {
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
			case "port":
				if index+1 < len(args) && isValidValue(args[index+1]) {
					argsMap[argument] = args[index+1]
					continue
				}
			case "replicaof":
				if index+1 < len(args) && isValidValue(args[index+1]) {
					argsMap["role"] = "slave"
					hostPort := args[index+1]
					hostPortSplit := strings.Split(hostPort, " ")
					if len(hostPortSplit) != 2 {
						panic("Invalid host and port")
					}
					argsMap["master"] = hostPortSplit[0]
					argsMap["master_port"] = hostPortSplit[1]
					continue
				}
			}
		}
		panic(fmt.Sprintf("Invalid argument %s", argument))
	}
}

func checkArguments(argsMap map[string]interface{}) {
	if _, ok := argsMap["dir"]; !ok {
		argsMap["dir"] = "."
	}
	if _, ok := argsMap["dbfilename"]; !ok {
		argsMap["dbfilename"] = "dump.rdb"
	}
	if _, ok := argsMap["port"]; !ok {
		argsMap["port"] = "6379"
	}
	if _, ok := argsMap["role"]; !ok {
		argsMap["role"] = "master"
	}
}

func isValidArg(arg string) bool {
	return strings.HasPrefix(arg, "--")
}

func isValidValue(value string) bool {
	return !strings.HasPrefix(value, "--") && (strings.TrimSpace(value) != "")
}
