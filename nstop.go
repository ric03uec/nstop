package main
import (
	"os"
	"log"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/ric03uec/nstop/arguments"
	"github.com/ric03uec/nstop/supervisor"
)

var NSTOP_CONFIG_FILENAME = ".nstopcfg.json"

func bootApplication(c *cli.Context) {
	// if file option set, then use the filename to load the config
	// else no config provided, use default
	log.Printf("Booting application")
	fileName := NSTOP_CONFIG_FILENAME
	var started bool
	var config *arguments.Config
	if c.IsSet("file"){
		fileName = c.GlobalString("file")
		config = arguments.Initialize(fileName)
	} else if len(c.Args()) > 0 {
		command_name := c.Args()[0]
		//TODO: get everything in the array and not just the first word, command
		// can be nstop node app.js -blah -blu -ble
		fmt.Printf("Booting application %s \n", command_name)
		//started, err := supervisor.Boot(config)
	}
	started, err := supervisor.Boot(config)
	if err == nil && started == true {
		log.Printf("Supervisor started successfully... ")
	} else {
		log.Printf("Error while booting logger: %v", err)
		os.Exit(1)
	}
	//logger.boot
	//watcher.boot
}

func main() {
	// use go-flags or getopt package for parsing flags
	// use channels
	// DONT pass arrays, pass slices
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

