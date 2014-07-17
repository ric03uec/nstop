package main
import (
	"os"
	"log"
	"github.com/codegangsta/cli"
	"github.com/ric03uec/nstop/arguments"
	"github.com/ric03uec/nstop/supervisor"
)

var NSTOP_CONFIG_FILENAME = ".nstopcfg.json"

func bootApplication(c *cli.Context) {
	log.Printf("Booting application")
	fileName := NSTOP_CONFIG_FILENAME
	if c.IsSet("file"){
		fileName = c.GlobalString("file")
	}
	config := arguments.Initialize(fileName)
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

