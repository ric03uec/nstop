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
	log.Printf("%v --- %v", supervisorConfig, err)
	if isValidConfig, err := validateConfig(supervisorConfig.Values); isValidConfig == false{
		log.Printf("Invalid configuation in config file")
		return false, err
	}
	configEntry, _ := supervisorConfig.GetConfigValue("exec")
	proc := NewProc(fmt.Sprintf("%s", configEntry.Value))
	proc.exec()
	log.Printf(fmt.Sprintf("%v", proc))

	return true, nil
}
