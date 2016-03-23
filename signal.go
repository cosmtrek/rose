package main

import (
	"os"
	"os/signal"
	"syscall"
)

func guardSignal() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-sigs
		done <- true
	}()
	<-done
	info.Println("Rose is exiting, bye!")
	pidFile := config.PidFile
	if _, err := os.Stat(pidFile); err == nil {
		os.Remove(pidFile)
	}
	os.Exit(0)
}
