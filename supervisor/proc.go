package supervisor

import (
	"time"
	"fmt"
	"syscall"
	//"bytes"
	"log"
	"os"
	"strings"
	"os/exec"
	"os/signal"
)

type Proc struct {
	cmd *exec.Cmd
	signalChannel chan os.Signal
	command string
	exitCode int16
	running bool
	restartCount uint16
	maxRestartCount uint16
	shell string
	waitTime uint16 //seconds
	procError string
}

func NewProc(command string) *Proc {
	proc := new(Proc)
	proc.shell = "/bin/bash"
	proc.waitTime = 5
	proc.maxRestartCount = 3
	proc.restartCount = 0
	proc.exitCode = -1
	proc.command = command
	proc.running = false
	return proc
}

func (proc *Proc) AddHandlers() {
	// add handlers for 
	// kill (kill everything)
	// term/int (gracefully kill)
	// usr/usr2 (restart, reload files)
}

func (proc *Proc) String() string {
	if proc.running == true {
		return fmt.Sprintf("Command: %s, Running: %t, Restart count: %d",
			proc.command, proc.running, proc.restartCount)
	} else {
		return fmt.Sprintf("Command: %s, Running: %t, Exit code: %d, Restart count: %d",
			proc.command, proc.running, proc.exitCode, proc.restartCount)
	}
}

func (proc *Proc) Stop() (safeStop bool, err error) {
	fmt.Printf("supervisor killing child process\n")
	proc.cmd.Process.Signal(syscall.SIGKILL)
	close(proc.signalChannel)
	return true, nil
}

func (proc *Proc) AddSignalHandlers() {
	proc.signalChannel = make(chan os.Signal, 1)
	signal.Notify(proc.signalChannel, os.Interrupt)
	//signal.Notify(proc.signalChannel)
	go func(){
		for sig := range proc.signalChannel {
			fmt.Printf("received signal at handler \n")
			fmt.Printf("%s", sig)
			proc.Stop()
			//proc.Stop(), if signal = int or stop
			//proc.ReloadChildConfig(), if signal=usr1
		}
	}()
}

func (proc *Proc) RemoveSignalHandlers() {

}

func (proc *Proc) Start() (safeExit bool, err error) {
	log.Printf("supervisor proc running command: %s\n", proc.command)
	commandParts := strings.Fields(proc.command)
	commandString := commandParts[0]
	commandArgs := commandParts[1:len(commandParts)]
	proc.cmd = exec.Command(commandString, commandArgs...)

	//store these in logs
	proc.cmd.Stdout = os.Stdout
	proc.cmd.Stderr = os.Stderr

	log.Printf("--------- Starting command execution ----------- \n")
	cmdErr := proc.cmd.Start()
	proc.AddSignalHandlers()
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}

	done := make(chan error, 1)
	go func() {
		done <- proc.cmd.Wait()
	}()

	select {
		case <- time.After(1200 * time.Second):
			fmt.Printf("killing process\n")
			proc.Stop()
			<-done
		case exitCode := <-done:
			if exitCode != nil {
				proc.exitCode = -1
				proc.procError = fmt.Sprintf("%v", exitCode)
				log.Printf("%v", proc.procError)
			} else {
				proc.exitCode = 0
			}
			log.Printf("--------- Command exited ----------- \n")
	}
	return true, nil
}
