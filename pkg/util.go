package pkg

import "strings"

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
