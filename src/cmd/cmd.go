package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"../log"
)

func exe_cmd(command string) bool {
	parts := strings.Fields(command)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(parts[0], parts[1])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		log.Error("[%v] - [%v]", err, stderr.String())
		return false
	}
	log.Info("Result: [%v]", out.String())
	return true
}

func Exec(command string) bool {
	return exe_cmd(command)
}
