package supervisor
import (
	"fmt"
	"errors"
	"log"
	"github.com/ric03uec/nstop/arguments"
)

var VALID_PARAMS = map[string]bool{
	"exec"		: true,
	"restartCount"	: true,
	"minWait"	: false,
	"numProcs"	: false,
}

var PARAM_DEFAULT = map[string]interface{}{
	"restartCount"	: 3,
	"minWait"	: 5,
	"numProcs"	: 1,
}

func GetDefaultConfig(exec string) []arguments.ModuleConfig {
	modulesConfigSlice := []arguments.ModuleConfig{}

	// adding default supervisor config
	defaultConfig := arguments.NewModuleConfig("supervisor")
	configExecEntry := arguments.NewConfigEntry("exec", exec)
	defaultConfig.Values = append(defaultConfig.Values, *configExecEntry)
	for key, _:= range VALID_PARAMS {
		if key == "exec" {
			continue
		}
		configValue := PARAM_DEFAULT[key]
		configEntry := arguments.NewConfigEntry(key, configValue)
		defaultConfig.Values = append(defaultConfig.Values, *configEntry)
	}

	//TODO: add default watcher/logger config
	modulesConfigSlice = append(modulesConfigSlice, *defaultConfig)

	return modulesConfigSlice
}

func isValidParamName(paramName string) bool {
	for key ,_ := range VALID_PARAMS {
		if paramName == key {
			return true
		}
	}
	return false
}

func validateConfig(configEntries []arguments.ConfigEntry) (isValid bool, err error) {
	isValid = true
	var validationError error
	for _, configEntry := range configEntries {
		if isValidParamName(configEntry.Key) == true {
			if VALID_PARAMS[configEntry.Key] == true && configEntry.Value == nil {
				isValid = false
				validationError = errors.New(fmt.Sprintf("Mandatory parameter not present: %s\n", configEntry.Key))
				break
			}
		} else {
			isValid = false
			validationError = errors.New(fmt.Sprintf("Invalid Key in params: %s\n", configEntry.Key))
			break
		}
	}
	return isValid, validationError
}

func Boot(config []arguments.ModuleConfig)(started bool, err error) {
	log.Printf("%v", config)
	supervisorConfig, err := arguments.GetModuleConfig(config, "supervisor")
	if isValidConfig, err := validateConfig(supervisorConfig.Values); isValidConfig == false{
		log.Printf("Invalid configuation in config file")
		return false, err
	}
	configEntry, _ := supervisorConfig.GetConfigValue("exec")
	proc := NewProc(fmt.Sprintf("%s", configEntry.Value))
	proc.Start()
	log.Printf(fmt.Sprintf("%v", proc))

	return true, nil
}
