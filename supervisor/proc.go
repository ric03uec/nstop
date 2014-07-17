package supervisor

import (
	//"fmt"
	//"bytes"
	"log"
	"os"
	"strings"
	"os/exec"
)

type Proc struct {
	command string
	exitCode int16
	running bool
	restartCount uint16
	shell string
	waitTime uint16 //seconds
}

func NewProc(command string) *Proc {
	proc := new(Proc)
	proc.shell = "/bin/bash"
	proc.waitTime = 3
	proc.restartCount = 0
	proc.exitCode = -1
	proc.command = command
	proc.running = false
	return proc
}


func (proc *Proc) exec() (safeExit bool, err error){
	log.Printf("supervisor proc running command: %s\n", proc.command)
	commandParts := strings.Fields(proc.command)
	commandString := commandParts[0]
	commandArgs := commandParts[1:len(commandParts)]
	cmd := exec.Command(commandString, commandArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("--------- Starting command execution ----------- \n")
	cmdErr := cmd.Start()
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
	err = cmd.Wait()
	log.Printf("--------- Command exited ----------- \n")
	//log.Printf("cmd finished %s", out.String())

	return true, nil
}
