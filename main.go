package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	args := os.Args[0:]

	if len(args) < 1 {
		log.Fatalf("Usage: %v cmd [ args ... ]\n", args[0])
	}

	program := args[1]
	programArgs := args[2:]

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	fmt.Printf("Running PID: %v\n", os.Getpid())
	go func() {
		for received_signal := range signals {
			switch received_signal {
			case syscall.SIGTERM:
			case syscall.SIGINT:
				fmt.Printf("received %v exiting\n", received_signal)
				done <- true
				break
			case syscall.SIGHUP:
				fmt.Printf("received %v running command %v args: %v \n", received_signal, program, programArgs)
				cmd := exec.Command(program, programArgs...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
				break
			}
		}
	}()
	<-done
}
