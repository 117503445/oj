package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func getStringWithLineLimit(text string, lineLimit int) string {
	lines := strings.Split(text, "\n")
	if len(lines) < lineLimit {
		return text
	}
	lines = lines[:lineLimit-1]
	lines = append(lines, "...") // todo
	text = strings.Join(lines, "\n")
	return text
}

func execBuild(exePath string) {
	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	cmd := exec.CommandContext(ctx, exePath)

	cmd.Stdin = strings.NewReader("1 2")

	out, _ := cmd.Output()

	text := getStringWithLineLimit(string(out), 5)

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Command timed out")
	}

	fmt.Println(text)
}

func main() {
	execBuild("./assets/success/build/main.exe")
	execBuild("./assets/timeout/build/main.exe")
	execBuild("./assets/output-limit-exceeded/build/main.exe")
}
