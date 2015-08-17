package main

import (
	"log"
	"strconv"
	"strings"
)

type CommandFunc func()

type commandContext struct {
	name     string
	enabled  bool
	commands map[string]CommandFunc
}

type Command struct {
	Name string
	Cmd  CommandFunc
}

var commandContexts []commandContext
var argumentStack [][]string

func RegisterCommands(context string, commands []Command) {
	contextIdx := -1
	for i := 0; i < len(commandContexts); i++ {
		if commandContexts[i].name == context {
			contextIdx = i
			break
		}
	}

	if contextIdx == -1 {
		commandContexts = append(commandContexts, commandContext{name: context, enabled: false, commands: make(map[string]CommandFunc)})
		contextIdx = len(commandContexts) - 1
	}

	for _, cmd := range commands {
		commandContexts[contextIdx].commands[cmd.Name] = cmd.Cmd
	}
}

func getCommand(name string) (CommandFunc, bool) {
	for i := len(commandContexts) - 1; i >= 0; i-- {
		cmd, exists := commandContexts[i].commands[name]
		if exists {
			return cmd, true
		}
	}

	return nil, false
}

func EnterCommandMode() {
	EnableInputMode(":", "", ExecLine)
}

func GetArg(idx int) (string, bool) {
	top := len(argumentStack) - 1
	if top >= 0 {
		if idx < len(argumentStack[top]) {
			return argumentStack[top][idx], true
		}
	}
	return "", false
}

func GetIntArg(idx int) (int, bool) {
	strValue, exists := GetArg(idx)
	if !exists {
		return 0, false
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		log.Printf("Error in command %s - Expected integer for argument %d but got '%s'", GetArgs()[0], idx, strValue)
		return 0, false
	}

	return intValue, true
}

func GetArgs() []string {
	if len(argumentStack) > 0 {
		return argumentStack[len(argumentStack)-1]
	}
	return nil
}

func GetArgCount() int {
	if len(argumentStack) > 0 {
		return len(argumentStack[len(argumentStack)-1])
	}
	return 0
}

func pushArgs(args []string) {
	argumentStack = append(argumentStack, args)
}

func popArgs() {
	argumentStack = argumentStack[:len(argumentStack)-1]
}

func ExecLine(line string) {
	pushArgs(strings.Fields(line))

	args := GetArgs()
	if len(args) >= 1 {
		cmd, exists := getCommand(args[0])
		if exists {
			cmd()
		}
	}

	popArgs()
}
