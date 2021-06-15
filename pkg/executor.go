package pkg

import (
	"context"
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

func ExecRun(outChan chan string, exePath string, inputPath string) {

	output := ""

	input, _ := ioutil.ReadFile(inputPath)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, exePath)

	cmd.Stdin = strings.NewReader(string(input))

	out, _ := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		output += "Command timed out\n"
	}

	text := GetStringWithLineLimit(string(out), 12)
	output += text + "\n"

	text = GetStringWithLineLimit(string(out), 500)
	err := ioutil.WriteFile(GetFileNameWithoutExt(inputPath)+".out", []byte(text), 0644)
	if err != nil {
		output += err.Error() + "\n"
	}

	outChan <- output
}
