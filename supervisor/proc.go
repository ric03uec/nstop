package supervisor

import (
	"fmt"
	"log"
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
	fmt.Printf("proc running command %s", proc.command)
	// run exec.CombinedOutput

	fmt.Sprintf("%s",strings.Split(proc.command, " "))
	cmd := exec.Command(proc.command)
	cmdErr := cmd.Start()
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
	log.Printf("waiting for cmd")
	err = cmd.Wait()
	log.Printf("cmd finished")

	return true, nil
}
