package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func exe_cmd(command string) bool {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1])

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return false
	}
	fmt.Println("Result: " + out.String())
	return true
}

func Exec(command string) bool {
	return exe_cmd(command)
}
