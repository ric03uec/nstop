package arguments
import (
	"os"
	"log"
)

func Initialize(fileName string)([]ModuleConfig) {
	//TODO: provide full path to the fiLe, easier for debuggin
	log.Printf("Reading configuration from file : %s\n", fileName)
	readConfig, err := NewConfig(fileName)
	if err != nil {
		log.Printf("Could not initialize configuration, exiting...")
		os.Exit(1)
	}
	return readConfig.parsedConfig
}
