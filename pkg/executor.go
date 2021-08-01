package pkg

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

// ExecBuild 执行编译，成功返回 "", 失败返回错误字符串
func ExecBuild(sourcePath string, language string) (success bool, output string) {
	key := fmt.Sprintf("languages.%v.build", language)
	command := BuildCommand(viper.GetString(key), sourcePath)
	//glog.Line().Debug(command)
	name, args := SplitCommand(command)
	cmd := exec.Command(name, args...)

	out, err := cmd.Output()
	if err != nil {
		//glog.Line().Error(err)
	}

	return err == nil, string(out)
}

func ExecRun(outChan chan string, sourcePath string, language string, inputPath string) {

	output := ""

	input, _ := ioutil.ReadFile(inputPath)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	key := fmt.Sprintf("languages.%v.run", language)
	command := BuildCommand(viper.GetString(key), sourcePath)

	cmd := exec.CommandContext(ctx, command)

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
