package conf

import (
	"encoding/json"
	"fmt"
	"github.com/go-yaml/yaml"
	"os"
	"regexp"
	"strconv"
)

// Local settings from config/{environment}.yaml
var Local Environment

// Global settings from config/global.yaml
var Global Environment

type Environment struct {
	Name string
}

// [environment][setting][value]
var configData map[string]map[string]interface{}
var re *regexp.Regexp

func init() {
	configData = make(map[string]map[string]interface{})
	re = regexp.MustCompile("^\\s*([\\w-]*)\\s*:\\s*(.*)\\s*")
	Global.Name = "global"
	if len(os.Args) > 1 {
		Local.Name = os.Args[1]
	} else {
		panic("Please run app with environment -> ./app environment")
	}
}

// GetEnv Return current environment, dev is default
func GetEnv() string {
	return Local.Name
}

// Get setting as string
func (e Environment) Get(setting string) (result string) {
	environmentMap := fetchenvironment(e)
	val, _ := environmentMap[setting]

	parse(val, &result)
	return
}

// GetUint get setting as uint64
func (e Environment) GetUint(setting string) uint64 {
	var result uint64

	environmentMap := fetchenvironment(e)
	val, _ := environmentMap[setting]
	parse(val, &result)

	return result
}

// GetInt get setting as int64
func (e Environment) GetInt(setting string) int64 {
	var result int64

	environmentMap := fetchenvironment(e)
	val, _ := environmentMap[setting]
	parse(val, &result)

	return result
}

// GetFloat get setting as float64
func (e Environment) GetFloat(setting string) float64 {
	var strVal string

	environmentMap := fetchenvironment(e)
	val, _ := environmentMap[setting]
	parse(val, &strVal)

	parsedVal, _ := strconv.ParseFloat(strVal, 64)
	return parsedVal
}

// GetBool get setting as boolean
func (e Environment) GetBool(setting string) bool {
	var strVal string

	environmentMap := fetchenvironment(e)
	val, _ := environmentMap[setting]
	parse(val, &strVal)

	parsedVal, _ := strconv.ParseBool(strVal)
	return parsedVal
}

// GetSlice get setting as slice of strings
func (e Environment) GetSlice(setting string) (result []string) {
	environmentMap := fetchenvironment(e)
	val, _ := environmentMap[setting]

	parse(val, &result)
	return
}

func fetchenvironment(e Environment) map[string]interface{} {
	environmentMap, ok := configData[e.Name]
	// singleton
	if !ok {
		importSettingsFromFile(e.Name)
		environmentMap, _ = configData[e.Name]
	}
	return environmentMap
}

func importSettingsFromFile(environment string) {
	config := make(map[string]interface{})
	file, err := os.ReadFile("config/" + environment + ".yaml")
	if err != nil {
		panic(fmt.Sprintf("Open config file fail: config/%s.yaml. Please run application as ./app [dev] ", environment))
		return
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(fmt.Sprintf("Parse config file fail: config/%s.yaml %s", environment, err.Error()))
		return
	}
	configData[environment] = config
}

func parse(in interface{}, out interface{}) {
	bytes, _ := json.Marshal(in)
	_ = json.Unmarshal(bytes, &out)
}
