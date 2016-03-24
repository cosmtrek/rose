package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
)

type infoLogging bool
type errorLogging bool
type debugLogging bool

var (
	info  infoLogging
	errl  errorLogging
	debug debugLogging

	logFile   io.Writer
	logFormat = log.LstdFlags

	infoLog  = log.New(os.Stdout, color.WhiteString("INFO: "), logFormat)
	errorLog = log.New(os.Stdout, color.RedString("ERROR: "), logFormat)
	debugLog = log.New(os.Stdout, color.BlueString("DEBUG: "), logFormat)
)

func initLog() {
	debugFlag := false
	if config.ServerEnv == "development" {
		debugFlag = true
	}

	flag.BoolVar((*bool)(&info), "info log", true, "show basic info")
	flag.BoolVar((*bool)(&errl), "error log", true, "show error info")
	flag.BoolVar((*bool)(&debug), "debug log", debugFlag, "show debug info")

	writePid()
}

func writePid() {
	pid := os.Getpid()
	pidFile := config.PidFile
	file, err := os.OpenFile(pidFile, os.O_WRONLY, os.ModeAppend)
	if err == nil {
		if _, err := file.WriteString(strconv.Itoa(pid)); err != nil {
			errl.Printf("Failed to write pid to file, error: %s", err)
		}
	} else {
		file, err = os.Create(pidFile)
		if err == nil {
			if _, err := file.WriteString(strconv.Itoa(pid)); err != nil {
				errl.Printf("Failed to write pid to file, error: %s", err)
			}
		}
	}
	file.Close()
}

func (l infoLogging) Printf(format string, args ...interface{}) {
	if l {
		infoLog.Printf(format, args...)
	}
}

func (l infoLogging) Println(args ...interface{}) {
	if l {
		infoLog.Println(args...)
	}
}

func (e errorLogging) Printf(format string, args ...interface{}) {
	if e {
		errorLog.Printf(format, args...)
	}
}

func (e errorLogging) Println(args ...interface{}) {
	if e {
		errorLog.Println(args...)
	}
}

func (d debugLogging) Printf(format string, args ...interface{}) {
	if d {
		debugLog.Printf(format, args...)
	}
}

func (d debugLogging) Println(args ...interface{}) {
	if d {
		debugLog.Println(args...)
	}
}

func Fatal(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}
