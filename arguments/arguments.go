package arguments
import (
	"os"
	"log"
)

func Initialize(fileName string)(config *Config) {
	log.Printf("Reading configuration from file : %s\n", fileName)
	readConfig, err := NewConfig(fileName)
	if err != nil {
		log.Printf("Could not initialize configuration, exiting...")
		os.Exit(1)
	}
	return readConfig
}
