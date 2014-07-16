package arguments
import (
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
)

type ConfigEntry struct {
	key string
	value interface{}
}

func NewConfigEntry(key string, value interface{}) *ConfigEntry {
	configEntry := new(ConfigEntry)
	configEntry.key = key
	configEntry.value = value
	return configEntry
}

type ModuleConfig struct {
	moduleName string
	values []ConfigEntry
}

func NewModuleConfig(moduleName string) *ModuleConfig {
	moduleConfig := new(ModuleConfig)
	moduleConfig.moduleName = moduleName
	moduleConfig.values = []ConfigEntry{}
	return moduleConfig
}

type Config struct {
	path string
	parsedConfig []ModuleConfig
}

func NewConfig(fileName string) *Config {
	config := new(Config)
	config.path = fileName
	config.parsedConfig = []ModuleConfig{}
	return config
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

func (config Config) ModuleConfig(moduleName string) *ModuleConfig {
	for _, module := range config.parsedConfig {
		if module.moduleName == moduleName {
			return &module
		}
	}
	return nil
}

func (config Config) FromFile() []ModuleConfig {
	content, err := ioutil.ReadFile(config.path)
	if err!=nil{
		fmt.Printf("Error", err)
	}

	var parsedJson map[string]interface{}
	json.Unmarshal([]byte(content), &parsedJson)

	for key, value := range parsedJson {
		if IsValidModule(key) {
			moduleConfig := NewModuleConfig(key)
			parsedModuleConfig := value.(map[string]interface{})
			for moduleKey, moduleVal := range parsedModuleConfig {
				configEntry := NewConfigEntry(moduleKey, moduleVal)
				moduleConfig.values = append(moduleConfig.values, *configEntry)
			}
			config.parsedConfig = append(config.parsedConfig , *moduleConfig)
		}
	}
	log.Printf("Parsed config from file %s\n", config.parsedConfig)
	return config.parsedConfig
}
