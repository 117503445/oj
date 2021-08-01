package pkg

import "strings"

// BuildCommand 替换命令中的占位符
func BuildCommand(raw string, codePath string) (command string) {
	binPath := GetBinFileName(codePath)
	raw = strings.ReplaceAll(raw, "${codePath}", codePath)
	raw = strings.ReplaceAll(raw, "${binPath}", binPath)
	return raw
}

func SplitCommand(command string) (name string, args []string) {
	words := strings.Split(command, " ")
	return words[0], words[1:]
}
