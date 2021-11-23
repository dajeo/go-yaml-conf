package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"os"
	"regexp"
	"strconv"
)

var configData map[string]map[string]string
var re *regexp.Regexp
var Local Environment
var Global Environment

type Environment struct {
	Name string
}

func init() {
	configData = make(map[string]map[string]string)
	re = regexp.MustCompile("^\\s*([\\w-]*)\\s*:\\s*(.*)\\s*")
	Global.Name = "global"
	if len(os.Args) > 1 {
		Local.Name = os.Args[1]
	} else {
		panic("Please run app with environment -> ./app environment")
	}
}

// Return current environment,  dev is default
func GetEnv() string {
	return Local.Name
}

func (e Environment) Get(setting string) string {
	environmentMap := fetchenvironment(e)
	val, _ := environmentMap[setting]
	return val
}

func (e Environment) GetUint(setting string) uint64 {
	environmentMap := fetchenvironment(e)
	val, _ := environmentMap[setting]
	parsedVal, _ := strconv.ParseUint(val, 10, 64)
	return parsedVal
}

func (e Environment) GetInt(setting string) int64 {
	environmentMap := fetchenvironment(e)
	val, _ := environmentMap[setting]
	parsedVal, _ := strconv.ParseInt(val, 10, 64)
	return parsedVal
}

func (e Environment) GetFloat(setting string) float64 {
	environmentMap := fetchenvironment(e)
	val, _ := environmentMap[setting]
	parsedVal, _ := strconv.ParseFloat(val, 64)
	return parsedVal
}

func (e Environment) GetBool(setting string) bool {
	environmentMap := fetchenvironment(e)
	val, _ := environmentMap[setting]
	parsedVal, _ := strconv.ParseBool(val)
	return parsedVal
}

//func (e Environment) GetSlice(setting string) []string {
//	environmentMap := fetchenvironment(e)
//	val, _ := environmentMap[setting]
//	parsedVal, _ := strconv.P(val)
//	return parsedVal
//}

func fetchenvironment(e Environment) map[string]string {
	environmentMap, ok := configData[e.Name]
	// singleton
	if !ok {
		importSettingsFromFile(e.Name)
		environmentMap, _ = configData[e.Name]
	}
	return environmentMap
}

func importSettingsFromFile(environment string) {
	configData[environment] = make(map[string]string)
	file, err := os.ReadFile("config/" + environment + ".yml")
	if err != nil {
		panic("Open config file fail: config/" + environment + ".yml. Please run application as ./app [dev] ")
		return
	}
	err = yaml.Unmarshal(file, &configData)
	if err != nil {
		panic(fmt.Sprintf("Parse config file fail: config/"+environment+".yml %s", err.Error()))
		return
	}
}
