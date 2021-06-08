package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"oj/pkg"
	"os"
	"time"
)

func main() {
	sourcePath := "main.cpp" // todo auto search

	var lastSourceContent []byte

	for true {
		sourceContent, _ := ioutil.ReadFile(sourcePath)
		if !bytes.Equal(lastSourceContent, sourceContent) {
			// file change

			pkg.Clear()

			isBuildSuccess := pkg.ExecBuild(sourcePath)
			if !isBuildSuccess {
				fmt.Println("build fail")
				continue
			}

			inputPathArray := []string{"1.txt", "2.txt"} // todo load from disk

			for _, inputPath := range inputPathArray {
				fmt.Printf("--- %s ---\n", inputPath)
				input, _ := ioutil.ReadFile(inputPath)
				pkg.ExecRun("1.exe", string(input))
				fmt.Printf("--- %s ---\n\n", inputPath)
			}
			os.Remove("1.exe")
		}

		lastSourceContent = sourceContent
		time.Sleep(time.Second)
	}

}
