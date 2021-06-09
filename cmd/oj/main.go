package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"oj/pkg"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	supportFileType := []string{".cpp"}

	var sourcePathArray []string

	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		for _, fileType := range supportFileType {
			if strings.HasSuffix(f.Name(), fileType) {
				sourcePathArray = append(sourcePathArray, f.Name())
			}
		}

	}

	sourcePath := ""

	switch len(sourcePathArray) {
	case 0:
		fmt.Println("No Source Code found in the dir.")
		return
	case 1:
		sourcePath = sourcePathArray[0]
	default:
		for i, file := range sourcePathArray {
			fmt.Printf("%v %s; ", i, file)
		}
		fmt.Println("\nplease input the index of source filename")

		index := 0
		fmt.Scanf("%d", &index)

		sourcePath = sourcePathArray[index]
		pkg.Clear()
	}

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

			inputPathArray, _ := filepath.Glob("*.in")

			for _, inputPath := range inputPathArray {
				fmt.Printf("--- %s ---\n", inputPath)
				pkg.ExecRun("./"+pkg.GetBinFileName(sourcePath), inputPath)
				fmt.Printf("--- %s ---\n\n", inputPath)
			}
			err := os.Remove(pkg.GetBinFileName(sourcePath))
			if err != nil {
				fmt.Println(err)
			}
		}

		lastSourceContent = sourceContent
		time.Sleep(time.Second)
	}

}
