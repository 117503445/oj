package main

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"oj/pkg"
	"os"
	"path/filepath"
	"strings"
)

func getSourcePath() string {
	mapLanguageExtension := make(map[string]string)

	for k, v := range viper.GetStringMap("languages") {
		mapLanguageExtension[k] = v.(map[string]interface{})["extension"].(string)
	}


	var sourcePathArray []string

	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		for language, extension := range mapLanguageExtension {
			if strings.HasSuffix(f.Name(), extension) {

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

func exec(ctx context.Context, sourcePath string) {
	pkg.Clear()
	var isBuildSuccess bool

	isBuildSuccess = pkg.ExecBuild(sourcePath)

	select {
	case <-ctx.Done():
		fmt.Println("exec cancel")
		return
	default:

	}

	if isBuildSuccess {
		inputPathArray, _ := filepath.Glob("*.in")

		channels := make([]chan string, 0)

		for _, inputPath := range inputPathArray {
			select {
			case <-ctx.Done():
				fmt.Println("exec cancel")
				return
			default:

			}
			channel := make(chan string)
			channels = append(channels, channel)

			go pkg.ExecRun(channel, "./"+pkg.GetBinFileName(sourcePath), inputPath)

		}

		for i, channel := range channels {
			inputPath := inputPathArray[i]
			output := ""
			output += fmt.Sprintf("--- %s ---\n", inputPath)
			output += <-channel
			output += fmt.Sprintf("--- %s ---\n\n", inputPath)
			fmt.Println(output)
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

	ctx, _ := context.WithCancel(context.Background())
	exec(ctx, sourcePath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	var lastCancel context.CancelFunc

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
					if lastCancel != nil {
						lastCancel()
					}

					ctx, cancel := context.WithCancel(context.Background())
					lastCancel = cancel
					exec(ctx, sourcePath)
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
