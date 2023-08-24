//go:build windows

// go build -ldflags -H=windowsgui

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/emmaly/winlaunch"
	"github.com/emmaly/winlaunch/config"
)

func main() {
	cfg := config.ReadConfig()
	program := config.ProgramConfig{}

	if len(cfg.Args) > 0 && cfg.Args[0] != "" {
		programName := cfg.Args[0]
		for _, p := range cfg.Programs {
			if strings.EqualFold(p.Name, programName) {
				program = p
				break
			}
		}
		if program.Name == "" {
			fmt.Println("Program not found in config file")
			os.Exit(1)
		}
	} else {
		fmt.Println("No program name provided")
		os.Exit(1)
	}

	err := winlaunch.Raise(program, cfg.Verbose)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
