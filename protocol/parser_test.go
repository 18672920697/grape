package protocol

import (
	"testing"
)

var testArray map[string]CommandData

func testJoint() {
	testArray := make(map[string]CommandData)

	req1 := "*3\r\n$3SET\r\n$3\r\nKEY\r\n$5\r\nVALUE\r\n"
	cmd1 := [3]string{"SET", "KEY", "VALUE"}
	testArray[req1] = CommandData{cmd1[:]}

	req2 := "+OK\r\n"
	cmd2 := [1]string{"OK"}
	testArray[req2] = CommandData{cmd2[:]}

	req3 := "-Deny heartbeat\r\n"
	cmd3 := [1]string{"Deny heartbeat"}
	testArray[req3] = CommandData{cmd3[:]}
}

func checkEqual(parse CommandData, standard CommandData, t *testing.T) {
	if len(parse.Args) != len(standard.Args) {
		t.Errorf("parse test error")
	}
	for i := range parse.Args {
		if parse.Args[i] != standard.Args[i] {
			t.Errorf("parse test error: %s - %s", parse.Args[i], standard.Args[i])
		}
	}
}

func TestParser(t *testing.T) {
	testJoint()

	for req, cmd := range testArray {
		cmdParse, _ := Parser(req)
		checkEqual(cmdParse, cmd, t)
	}
}
