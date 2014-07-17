package arguments
import (
	"fmt"
	"errors"
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	path string
	parsedConfig []ModuleConfig
}

type ModuleConfig struct {
	ModuleName string
	Values []ConfigEntry
}

type ConfigEntry struct {
	Key string
	Value interface{}
}

func NewConfigEntry(key string, value interface{}) *ConfigEntry {
	configEntry := new(ConfigEntry)
	configEntry.Key = key
	configEntry.Value = value
	return configEntry
}

func NewModuleConfig(moduleName string) *ModuleConfig {
	moduleConfig := new(ModuleConfig)
	moduleConfig.ModuleName = moduleName
	moduleConfig.Values = []ConfigEntry{}
	return moduleConfig
}

func NewConfig(fileName string) (config *Config, err error) {
	newConfig := new(Config)
	dir,_ := os.Getwd()
	newConfig.path = fmt.Sprintf("%s/%s", dir, fileName)
	newConfig.parsedConfig, err = LoadConfig(newConfig.path)
	if err != nil {
		return nil, err
	} else {
		return newConfig, nil
	}
}

func (config *Config) GetModuleConfig(moduleName string) (moduleConfig ModuleConfig, err error) {
	for _, config := range config.parsedConfig {
		if config.ModuleName == moduleName {
			return config, nil
		}
	}
	newConfig := NewModuleConfig(moduleName)
	return *newConfig, errors.New(fmt.Sprintf("Empty module config : %s", moduleName))
}

func (moduleConfig *ModuleConfig) GetConfigValue(key string) (configEntry ConfigEntry, err error){
	for _, moduleConfigEntry := range moduleConfig.Values {
		if moduleConfigEntry.Key == key {
			return moduleConfigEntry, nil
		}
	}
	newConfigEntry := NewConfigEntry("nil", nil)
	return *newConfigEntry, errors.New(fmt.Sprintf("Key %s not present in module config : %s", key, moduleConfig.ModuleName))
}

func LoadConfig(filePath string) (moduleConfig []ModuleConfig, err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil{
		fmt.Printf("Error while reading file: %v\n", err)
		return nil, err
	}
	modulesConfigSlice := []ModuleConfig{}

	var parsedJson map[string]interface{}
	json.Unmarshal([]byte(content), &parsedJson)

	for key, value := range parsedJson {
		if IsValidModule(key) {
			moduleConfig := NewModuleConfig(key)
			parsedModuleConfig := value.(map[string]interface{})
			for moduleKey, moduleVal := range parsedModuleConfig {
				configEntry := NewConfigEntry(moduleKey, moduleVal)
				moduleConfig.Values = append(moduleConfig.Values, *configEntry)
			}
			modulesConfigSlice = append(modulesConfigSlice, *moduleConfig)
		}
	}
	log.Printf("Parsed config from file %s\n", modulesConfigSlice)
	return modulesConfigSlice, nil
}

func IsValidModule(moduleName string) bool {
	validModules := []string{"supervisor", "watcher", "logger"}
	for _, module := range validModules {
		if moduleName == module {
			return true
		}
	}
	return false
}
