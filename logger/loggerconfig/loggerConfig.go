package loggerconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/viper"
)

//Init is a config intializer function
func Init() {
	jsonPath, _ := filepath.Abs("../logger/loggerconfig/loggerConfig.json")
	bodyBytes, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		fmt.Println("Couldn't read local configuration json file.", err)
	}

	parseConfiguration(bodyBytes)
}

func parseConfiguration(body []byte) {
	var localConfig loggerConfig
	err := json.Unmarshal(body, &localConfig)
	if err != nil {
		fmt.Println("Cannot parse logger configuration from json file, message: " + err.Error())
	}
	for key, value := range localConfig.Logger[0].Source {
		viper.Set(key, value)
	}
}

type loggerConfig struct {
	Logger []logger `json:"logger"`
}

type logger struct {
	Source map[string]interface{} `json:"source"`
}
