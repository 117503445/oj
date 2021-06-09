package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func GetStringWithLineLimit(text string, lineLimit int) string {
	lines := strings.Split(text, "\n")
	if len(lines) < lineLimit {
		return text
	}
	lines = lines[:lineLimit-1]
	lines = append(lines, "...") // todo
	text = strings.Join(lines, "\n")
	return text
}

func GetBinFileName(sourceFileName string) string {
	switch runtime.GOOS {
	case "windows":
		return sourceFileName + ".exe"
	default:
		return sourceFileName + ".bin"
	}
}

// https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	//todo mac
}

func Clear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic(fmt.Sprintf("Your platform [%s] is unsupported! I can't clear terminal screen :(", runtime.GOOS))
	}
}
