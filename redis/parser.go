// parser
package redis

import (
	"strings"
)

type Status int

type CommandData struct {
	Args []string
}

// Parse request
func Parser(request string) (CommandData, error) {
	if request[0] == '+' {
		return simpleString(request)
	} else if request[0] == '*' {
		return array(request)
	} else if request[0] == '-' {
		return errors(request)
	} else {
		return CommandData{}, nil
	}
}

func simpleString(request string) (CommandData, error) {
	var args []string
	split := strings.Split(request[1:], "\r\n")

	args = append(args, split[0])

	return CommandData{
		args,
	}, nil
}

func array(request string) (CommandData, error) {
	var args []string
	split := strings.Split(request, "\r\n")

	len := len(split)

	for index := 2; index < len; index = index + 2 {
		args = append(args, split[index])
	}
	return CommandData{
		args,
	}, nil
}

func errors(request string) (CommandData, error) {
	var args []string
	split := strings.Split(request[1:], "\r\n")

	args = append(args, split[0])

	return CommandData{
		args,
	}, nil
}
