package main

import (
	"fmt"
	"os"
	"os/exec"
)

var cmds = []string{
	"test",
	"-cover",
	"-coverprofile=coverage.out",
	"-v",
}

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1:]
		cmds = append(cmds, args...)
	}
	cmd := exec.Command("go", cmds...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}

	fmt.Print(string(output))
}
