package protocol

import (
	"strings"
	"strconv"
)

type CommandData struct {
	length int
	Command string
	Args []string
}


func Parser(request string) (CommandData, error) {
	split := strings.Split(request, "\r\n")
	len, _ := strconv.Atoi(split[0])
	return CommandData {
		len,
		split[1],
		split[2:],
	}, nil
}
