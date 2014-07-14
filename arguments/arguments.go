package arguments
import (
	"fmt"
	"log"
)

func Initialize(arguments []string) {
	log.Printf("App arguments %s\n", arguments)
	fmt.Println(arguments)
}
