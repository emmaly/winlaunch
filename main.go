package main

import (
	"fmt"
	"os"
	"os/exec"
)

var verboseOutput = false

func main() {
	config := readConfig()
	verboseOutput = config.Verbose

	windowRaised := raiseMatchedWindow(config.Program.TitleMatch)
	if windowRaised {
		if verboseOutput {
			fmt.Println("Window matched; raised.")
		}
		os.Exit(0)
	}

	if verboseOutput {
		fmt.Println("Window not matched; spawning new...")
	}
	launched := launchProgram(config.Program)

	if launched {
		if verboseOutput {
			fmt.Println("Launched.")
		}
		os.Exit(0)
	}

	if verboseOutput {
		fmt.Println("Failed to launch.")
	}
	os.Exit(1)
}

func launchProgram(program ProgramConfig) bool {
	cmd := exec.Command(program.Command, program.Args...)
	if program.StartIn != "" {
		cmd.Dir = program.StartIn
	}

	err := cmd.Start()
	return err == nil
}
