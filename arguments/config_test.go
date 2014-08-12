package arguments

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestInvalidConfigFile(t *testing.T) {
	_, err := NewConfig("invalidFileName")
	if err == nil {
		t.Error("config should be nil if invalid filename provided")
	} else {
		t.Log("No config created if invalid filename provided")
	}
}

func TestValidConfigFile(t *testing.T) {
	tempConfigFile := "configTest"
	dir,_ := os.Getwd()
	file, err := ioutil.TempFile(dir, tempConfigFile)
	if err != nil {
		t.Error("unable to create temporary config file")
	} else {
		fileData := []byte("{'supervisor': {}}")
		file.Write(fileData);
		defer file.Close()
		defer os.Remove(file.Name())
		_, configErr := NewConfig(tempConfigFile)
		if configErr == nil {
			t.Error("Unable to read from config file")
		} else {
			t.Log("Successfully read config file")
		}
	}
}

