package supervisor
import (
	"fmt"
	"errors"
	"log"
	"github.com/ric03uec/nstop/arguments"
)

var VALID_PARAMS = map[string]bool{
	"exec"		: true,
	"retryCount"	: true,
	"minWait"	: false,
	"numProcs"	: false,
}

var PARAM_DEFAULT = map[string]interface{}{
	"retryCount"	: 3,
	"minWait"	: 5,
	"numProcs"	: 1,
}

func getDefaultConfig(exec string) *arguments.ModuleConfig {
	defaultConfig := arguments.NewModuleConfig("supervisor")
	configExecEntry := arguments.NewConfigEntry("exec", exec)
	defaultConfig.Values = append(defaultConfig.Values, *configExecEntry)
	fmt.Printf("blahblah\n")

	return defaultConfig
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

func Boot(config *arguments.Config)(started bool, err error) {
	log.Printf("%v", config)
	supervisorConfig, err := config.GetModuleConfig("supervisor")
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
