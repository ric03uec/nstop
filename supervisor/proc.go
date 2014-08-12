package supervisor

import (
	"time"
	"fmt"
	"syscall"
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

func (proc *Proc) String() string {
	if proc.running == true {
		return fmt.Sprintf("Command: %s, Running: %t, Restart count: %d",
			proc.command, proc.running, proc.restartCount)
	} else {
		return fmt.Sprintf("Command: %s, Running: %t, Exit code: %d, Restart count: %d",
			proc.command, proc.running, proc.exitCode, proc.restartCount)
	}
}

func (proc *Proc) Stop(sig os.Signal) (safeStop bool, err error) {
	fmt.Printf("supervisor killing child process with signal: %s\n", sig)
	proc.cmd.Process.Signal(sig)
	signal.Stop(proc.signalChannel)
	return true, nil
}

func (proc *Proc) AddSignalHandlers() {
	proc.signalChannel = make(chan os.Signal, 1)
	signal.Notify(proc.signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go proc.StartSignalListener()
}

func (proc *Proc) StartSignalListener() {
	// filter signals for handling proc
	// usr/usr2 (restart, reload files)
	for sig := range proc.signalChannel {
		fmt.Printf("received signal at handler: %s\n", sig)
		proc.Stop(sig)
	}
}

func (proc *Proc) Start() (safeExit bool, err error) {
	log.Printf("supervisor running command: %s\n", proc.command)
	commandParts := strings.Fields(proc.command)
	commandString := commandParts[0]
	commandArgs := commandParts[1:len(commandParts)]

	proc.cmd = exec.Command(commandString, commandArgs...)

	proc.cmd.Stdout = os.Stdout
	proc.cmd.Stderr = os.Stderr

	proc.AddSignalHandlers()
	cmdErr := proc.cmd.Start()
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
			proc.Stop(os.Interrupt)
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
