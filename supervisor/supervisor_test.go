package supervisor

import (
	"testing"
	"fmt"
)

func TestDefaultExec(t *testing.T) {
	// tests whether the default config has the correct command or not
	fmt.Printf("test run failed")
	exec_cmd := "node app.js"
	supervisorConfig := getDefaultConfig(exec_cmd)
	configEntry, _ := supervisorConfig.GetConfigValue("exec")

	if configEntry.Value == exec_cmd {
		t.Log("success")
	} else {
		t.Error("default exec command does not match arguments")
	}
}
