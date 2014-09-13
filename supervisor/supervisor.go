package supervisor
import (
	"fmt"
	"errors"
	"log"
	"github.com/ric03uec/nstop/arguments"
)

var VALID_PARAMS_SUPERVISOR = map[string]bool{
	"exec"		: true,
	"restartCount"	: true,
	"minWait"	: false,
	"numProcs"	: false,
}

var PARAM_DEFAULT_SUPERVISOR = map[string]interface{}{
	"restartCount"	: 3,
	"minWait"	: 5,
	"numProcs"	: 1,
}

var VALID_PARAMS_WATCHER = map[string]bool{
	"enabled"	: true,
	"ignorePattern"	: false,
}

var PARAM_DEFAULT_WATCHER = map[string]interface{}{
	"enabled"	: true,
	"ignorePattern"	: []string{},
}

func GetDefaultConfig(exec string) []arguments.ModuleConfig {
	modulesConfigSlice := []arguments.ModuleConfig{}

	// adding default supervisor config
	defaultSupervisorConfig := arguments.NewModuleConfig("supervisor")
	configExecEntry := arguments.NewConfigEntry("exec", exec)
	defaultSupervisorConfig.Values = append(defaultSupervisorConfig.Values, *configExecEntry)
	for key, _:= range VALID_PARAMS_SUPERVISOR {
		if key == "exec" {
			continue
		}
		configValue := PARAM_DEFAULT_SUPERVISOR[key]
		configEntry := arguments.NewConfigEntry(key, configValue)
		defaultSupervisorConfig.Values = append(defaultSupervisorConfig.Values, *configEntry)
	}

	defaultWatcherConfig := arguments.NewModuleConfig("watcher")
	for key, _:= range VALID_PARAMS_WATCHER {
		configValue := PARAM_DEFAULT_WATCHER[key]
		configEntry := arguments.NewConfigEntry(key, configValue)
		defaultWatcherConfig.Values = append(defaultWatcherConfig.Values, *configEntry)
	}

	modulesConfigSlice = append(modulesConfigSlice, *defaultSupervisorConfig)
	modulesConfigSlice = append(modulesConfigSlice, *defaultWatcherConfig)

	return modulesConfigSlice
}

func isValidParamName(validParamMap map[string]bool, paramName string) bool {
	for key ,_ := range validParamMap {
		if paramName == key {
			return true
		}
	}
	return false
}

func validateConfig(validParamMap map[string]bool, configEntries []arguments.ConfigEntry) (isValid bool, err error) {
	isValid = true
	var validationError error
	for _, configEntry := range configEntries {
		fmt.Sprintf("%v", validParamMap)
		if isValidParamName(validParamMap, configEntry.Key) == true {
			if validParamMap[configEntry.Key] == true && configEntry.Value == nil {
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

// this will import the watcher config 
// proc will do the file watching and take decision based on whether the file has changed or notj
func Boot(config []arguments.ModuleConfig)(started bool, err error) {
	log.Printf("%v", config)
	supervisorConfig, err := arguments.GetModuleConfig(config, "supervisor")
	watcherConfig, err := arguments.GetModuleConfig(config, "watcher")

	if isValidConfig, err := validateConfig(VALID_PARAMS_SUPERVISOR, supervisorConfig.Values); isValidConfig == false{
		log.Printf("Invalid supervisor config")
		return false, err
	}
	configEntry, _ := supervisorConfig.GetConfigValue("exec")
	proc := NewProc(fmt.Sprintf("%s", configEntry.Value))
	proc.supervisorConfig = supervisorConfig

	if isValidConfig, err := validateConfig(VALID_PARAMS_WATCHER, watcherConfig.Values); isValidConfig == false{
		log.Printf("Invalid watcher config")
		return false, err
	}
	proc.watcherConfig = watcherConfig

	safeExit, procErr := proc.Start()
	if procErr == nil && safeExit == true {
		log.Printf(fmt.Sprintf("%v", proc))
		return true, nil
	} else {
		log.Printf(fmt.Sprintf("Error while starting command"))
		return safeExit, procErr
	}
}
