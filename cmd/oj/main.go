package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"oj/pkg"
	"os"
	"path/filepath"
	"strings"
)

func getSourcePath() string {
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
		return ""
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
	return sourcePath
}

func exec(sourcePath string) {
	pkg.Clear()

	isBuildSuccess := pkg.ExecBuild(sourcePath)
	if isBuildSuccess {
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
	} else {
		fmt.Println("build fail")
	}
}

func main() {

	sourcePath := getSourcePath()
	if sourcePath == "" {
		fmt.Println("No Source Code found in the dir.")
		return
	}

	exec(sourcePath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//fmt.Println(event.String())
				//fmt.Println(event.Name)
				//fmt.Println(sourcePath)
				//fmt.Println(strings.Contains(event.Name, sourcePath))
				//fmt.Println()
				if (strings.Contains(event.Name, sourcePath)) || strings.HasSuffix(event.Name, ".in") {
					exec(sourcePath)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
