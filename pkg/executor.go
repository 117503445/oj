package pkg

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func ExecBuild(sourcePath string) bool {
	// todo different build arg
	cmd := exec.Command("g++", sourcePath, "-o", "1.exe")

	err := cmd.Run()

	return err == nil
}

func ExecRun(exePath string) {
	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	cmd := exec.CommandContext(ctx, exePath)

	cmd.Stdin = strings.NewReader("1 2")

	out, _ := cmd.Output()

	text := GetStringWithLineLimit(string(out), 5)

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Command timed out")
	}

	fmt.Println(text)
}
