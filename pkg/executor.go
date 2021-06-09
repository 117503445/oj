package pkg

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

func ExecBuild(sourcePath string) bool {
	// todo different build arg
	cmd := exec.Command("g++", sourcePath, "-o", GetBinFileName(sourcePath))

	err := cmd.Run()

	return err == nil
}

func ExecRun(exePath string, inputPath string) {

	input, _ := ioutil.ReadFile(inputPath)

	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	cmd := exec.CommandContext(ctx, exePath)

	cmd.Stdin = strings.NewReader(string(input))

	out, _ := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Command timed out")
	}

	text := GetStringWithLineLimit(string(out), 8)
	fmt.Println(text)

	text = GetStringWithLineLimit(string(out), 100)
	err := ioutil.WriteFile(GetFileNameWithoutExt(inputPath)+".out", []byte(text), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
