package supervisor

import (
	"fmt"
	"syscall"
	"log"
	"os"
	"strings"
	"os/exec"
	"os/signal"
	"gopkg.in/fsnotify.v1"
	"github.com/ric03uec/nstop/arguments"
)

type Proc struct {
	cmd *exec.Cmd
	signalChannel chan os.Signal
	command string
	exitCode int
	running bool
	restartCount uint16
	maxRestartCount uint16
	shell string
	waitTime uint16 //seconds
	procError string
	supervisorConfig arguments.ModuleConfig
	watcherConfig arguments.ModuleConfig
	watcher *fsnotify.Watcher
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
	proc.watcher.Close()
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
		fmt.Printf("received signal at handler: %s, sending SIGKILL to process\n", sig)
		proc.Stop(syscall.SIGKILL)
	}
}

func (proc *Proc) InitWatcher() (success bool, err error) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Printf("Error while intializing watcher, exiting")
		return false, err
	}
	proc.watcher = watcher

	go func() {
		for {
			select {
			case ev := <- proc.watcher.Events:
				if ev.Op != 0 {
					log.Println("File changed: ", ev.Name)
					log.Println("OP: ", ev.Op)
					//??? how to stop this
					proc.Start()
				}
			case err := <- proc.watcher.Errors:
				if err != nil {
					log.Println("watcher err: ", err)
				}
			}
		}
	}()

	err = proc.watcher.Add(".")
	if err != nil {
		log.Printf("Error while watching directory")
		return false, err
	}
	return true, nil
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
	go func(){
		done <- proc.cmd.Wait()
	}()

	// add the watcher here
	proc.InitWatcher()
	//process waits here
	exitCode := <-done
	if exitCode != nil {
		// Type Assertion of exitCode with exec.ExitError struct
		var waitStatus syscall.WaitStatus
		exitError, ok := exitCode.(*exec.ExitError);
		if ok {
			// exitError.Sys() returns system specific exit info
			// Type Assertion of exitError.Sys() to sysCall.Waitstatus for Unix
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			exitStatus := waitStatus.ExitStatus()
			log.Printf(fmt.Sprintf("%v", waitStatus))
			if exitStatus == -1 {
				log.Printf("Process killed manually, aborting restart \n")
				proc.Stop(syscall.SIGKILL)
				return true, nil
			} else {
				proc.exitCode = exitStatus
				proc.procError = fmt.Sprintf("%v", exitCode)
				log.Printf("%v", proc.procError)
				if(proc.restartCount < proc.maxRestartCount){
					log.Printf("Restarting app, Restart count %d, Max restart count %d", proc.restartCount, proc.maxRestartCount)
					proc.restartCount += 1;
					proc.Start()
				}else{
					log.Printf("Max restart count (%d) reached, exiting", proc.maxRestartCount)
					log.Printf("--------- Command exited ----------- \n")
				}
			}
		}
	} else {
		// command exited successfully with return code = 0
		proc.exitCode = 0
		proc.restartCount += 1
		proc.Start()
	}
	return true, nil
}
