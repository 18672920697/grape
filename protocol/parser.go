// parser
package protocol

import (
	"strconv"
	"strings"
)

type Status int

type CommandData struct {
	length int
	Args   []string
}

func Parser(request string) (CommandData, error) {
	var args []string

	split := strings.Split(request, "\r\n")
	flag := split[0][1]
	len := len(split)

	cmd_len, _ := strconv.Atoi(string(flag))
	for index := 2; index < len; index = index + 2 {
		args = append(args, split[index])
	}
	return CommandData{
		cmd_len,
		args,
	}, nil
}
