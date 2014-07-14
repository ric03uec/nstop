package main
import (
	"os"
	"log"
	"github.com/codegangsta/cli"
	"github.com/ric03uec/nstop/arguments"
)

func bootApplication(c *cli.Context) {
	fileName := ".nstopcfg"
	if c.IsSet("file"){
		fileName = c.GlobalString("file")
	}
	log.Printf("Starting supervisor using config file : %s ", fileName)
	arguments.Initialize(fileName)
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

