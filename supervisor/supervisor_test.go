package supervisor

import (
	"testing"
	"fmt"
	"github.com/ric03uec/nstop/arguments"
)

func TestDefaultExec(t *testing.T) {
	// tests whether the default config has the correct command or not
	exec_cmd := "node app.js"
	defaultConfig := GetDefaultConfig(exec_cmd)
	supervisorConfig, _ := arguments.GetModuleConfig(defaultConfig, "supervisor")
	configEntry, _ := supervisorConfig.GetConfigValue("exec")

	if configEntry.Value == exec_cmd {
		t.Log("success")
	} else {
		t.Error("default exec command does not match arguments")
	}
}

func TestDefaultParamsCount(t *testing.T) {
	// test whether the default config has correct number of parameters
	exec_cmd := "node app.js"
	defaultConfig := GetDefaultConfig(exec_cmd)
	supervisorConfig, _ := arguments.GetModuleConfig(defaultConfig, "supervisor")

	configParamsCount := len(supervisorConfig.Values)
	validParamsCount := len(VALID_PARAMS)
	if configParamsCount == validParamsCount {
		fmt.Printf("%v\n", supervisorConfig)
		t.Log("success")
	} else {
		fmt.Printf("%v\n", supervisorConfig)
		//TODO: multiliine go statement
		t.Error(fmt.Sprintf("Invalid config entries count, \nValid: %d, Received %d\n", validParamsCount, configParamsCount))
	}
}

func TestDefaultParamsPresent(t *testing.T) {
	// test whether the default config has entries for all prameters

}
