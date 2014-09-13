package main
import (
	"os"
	"log"
	"fmt"
	"strings"
	"github.com/codegangsta/cli"
	"github.com/ric03uec/nstop/arguments"
	"github.com/ric03uec/nstop/supervisor"
)

var NSTOP_CONFIG_FILENAME = ".nstopcfg.json"

func bootApplication(c *cli.Context) {
	log.Printf("Booting application")
	fileName := NSTOP_CONFIG_FILENAME
	var config []arguments.ModuleConfig
	if c.IsSet("file"){
		fileName = c.GlobalString("file")
		config = arguments.Initialize(fileName)
	} else if len(c.Args()) > 0 {
		exec_cmd := fmt.Sprintf("%s", strings.Join(c.Args(), " "))
		config = supervisor.GetDefaultConfig(exec_cmd)
		fmt.Printf("Booting application %s \n", exec_cmd)
	}
	safeExit, err := supervisor.Boot(config)
	if err == nil && safeExit == true {
		log.Printf("Supervisor exited successfully... ")
	} else {
		log.Printf("Error while booting supervisor : %v", err)
		os.Exit(1)
	}
}

func main() {
	// use go-flags or getopt package for parsing flags
	app := cli.NewApp()
	app.Name = "!stop"
	app.Usage = "supervisor for docker applications"
	app.Action = bootApplication
	app.Version = "0.0.1"

	appFlags := []cli.Flag {
		cli.StringFlag{Name: "file, f", Usage: "configuration file name"},
	}
	app.Flags = appFlags
	app.Run(os.Args)
}

