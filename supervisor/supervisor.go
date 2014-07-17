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
		fmt.Printf("-------------------------\n")
		fmt.Printf("%v\n",configEntry.Key)
		fmt.Printf("%v\n",VALID_PARAMS[configEntry.Key])
		if isValidParamName(configEntry.Key) == true {
			if VALID_PARAMS[configEntry.Key] == true && configEntry.Value == nil {
				isValid = false
				validationError = errors.New("Mandatory parameter not present")
				break
			}
		} else {
			isValid = false
			validationError = errors.New("Incorrect Key in params")
			break
		}

	}
	fmt.Printf("%s\n",isValid)
	return isValid, validationError
}

func Boot(config *arguments.Config)(started bool, err error) {
	log.Printf("%v", config)
	supervisorConfig, err := config.GetModuleConfig("supervisor")
	log.Printf("%v --- %v", supervisorConfig, err)
	if isValidConfig, err := validateConfig(supervisorConfig); isValidConfig == false{
		log.Printf("Invalid Configuration file")
		return false, err
	}
	//exec command

	return true, nil
}

func PrintSth() {
	fmt.Printf("inside supervisor \n")
}
