package arguments
import (
	"fmt"
	"os"
	"log"
)

func Initialize(fileName string) {
	log.Printf("Reading configuration from file : %s\n", fileName)
	readConfig, err := NewConfig(fileName)
	if err != nil {
		log.Printf("Could not initialize configuration, exiting...")
		os.Exit(1)
	} else {
		fmt.Printf("%v", readConfig)
	}
}
