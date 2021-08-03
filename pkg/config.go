package pkg

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

const defaultConfJson = `{
    "languages": {
        "python": {
            "extension": "py",
            "run": "python ${codePath}"
        },
        "cpp": {
            "extension": "cpp",
            "build": "g++ ${codePath} -o ${binPath}",
            "run": "./${binPath}"
        }
    },
	"executor":{
		"outputLimit":{
			"terminal": 12,
			"file": 500
		}
	}
}`

func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name

	// TODO load config in file

	err := viper.ReadConfig(strings.NewReader(defaultConfJson))
	if err != nil {
		log.Panic(err)
		return
	}
}
