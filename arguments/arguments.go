package arguments
import (
	"log"
)

func Initialize(fileName string) {
	log.Printf("Reading configuration from file : %s\n", fileName)
	config := NewConfig(fileName)
	// constructor reads the file and inserts valid json or not
	//

	// config.FromFile()
	//	loads the config from provided json
	log.Printf(config.path)
	config.FromFile()
	// bubble read|parsing error to the caller
	// retu
}
