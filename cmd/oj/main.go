package main

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"oj/pkg"
	"os"
	"path/filepath"
	"strings"
)

// 通过用户输入,获得 源代码 的 路径 language
func getCodeInfo() (string, string) {
	mapLanguageExtension := make(map[string]string)

	for k, v := range viper.GetStringMap("languages") {
		mapLanguageExtension[k] = v.(map[string]interface{})["extension"].(string)
	}

	mapPathLanguage := make(map[string]string)

	files, _ := os.ReadDir("./")

	for _, f := range files {
		for language, extension := range mapLanguageExtension {
			if strings.HasSuffix(f.Name(), extension) {
				mapPathLanguage[f.Name()] = language
			}
		}
	}

	switch len(mapPathLanguage) {
	case 0:
		return "", ""
	case 1:
		for k, v := range mapPathLanguage {
			return k, v
		}

		panic("len(mapPathLanguage) == 1, but can't get items.")
	default:
		paths := make([]string, 0, len(mapPathLanguage))
		for k := range mapPathLanguage {
			paths = append(paths, k)
		}

		for i, path := range paths {
			fmt.Printf("%v %s; ", i, path)
		}

		fmt.Println("\nplease input the index of source filename")

		index := 0
		_, err := fmt.Scanf("%d", &index)
		if err != nil {
			panic(err)
		}

		pkg.Clear()

		if index < 0 || index >= len(paths) {
			panic("illegal index")
		} else {
			path := paths[index]
			language := mapPathLanguage[path]
			return path, language
		}
	}

}

func exec(ctx context.Context, sourcePath string, language string) {
	pkg.Clear()

	key := fmt.Sprintf("languages.%v.build", language)
	needCompile := viper.Get(key) != nil // 存在 languages.compile 时，才进行 compile

	if needCompile {
		// Build
		buildSuccess, buildOutput := pkg.ExecBuild(sourcePath, language)

		select {
		case <-ctx.Done():
			fmt.Println("exec cancel")
			return
		default:

		}
		if !buildSuccess {
			fmt.Println("Build Failed\n" + buildOutput)
			return
		}
	}

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

		go pkg.ExecRun(channel, sourcePath, language, inputPath)

	}

	for i, channel := range channels {
		inputPath := inputPathArray[i]
		output := ""
		output += fmt.Sprintf("--- %s ---\n", inputPath)
		output += <-channel
		output += fmt.Sprintf("--- %s ---\n\n", inputPath)
		fmt.Println(output)
	}

	if needCompile {
		// 删除可执行文件
		err := os.Remove(pkg.GetBinFileName(sourcePath))
		if err != nil {
			fmt.Println(err)
		}
	}

}

func main() {
	pflag.Int("executor.outputLimit.terminal", 12, "output lines limit in terminal")
	pflag.Int("executor.outputLimit.file", 500, "output lines limit in file")

	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		fmt.Println(err)
		return
	}

	sourcePath, language := getCodeInfo()
	if sourcePath == "" {
		fmt.Println("No Source Code found in the dir.")
		return
	}

	ctx, _ := context.WithCancel(context.Background())
	exec(ctx, sourcePath, language)

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
				if strings.HasSuffix(event.Name, sourcePath) || strings.HasSuffix(event.Name, ".in") {
					if lastCancel != nil {
						lastCancel()
					}

					ctx, cancel := context.WithCancel(context.Background())
					lastCancel = cancel
					exec(ctx, sourcePath, language)
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
