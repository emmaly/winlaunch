package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/emmaly/winlaunch/nillaBool"
)

type ProgramConfig struct {
	Name       string   `json:"name"`
	TitleMatch string   `json:"titleMatch"`
	Command    string   `json:"command"`
	Args       []string `json:"args"`
	StartIn    string   `json:"startIn"`
}

type Config struct {
	Verbose  bool            `json:"verbose"`
	Programs []ProgramConfig `json:"programs"`
	Program  ProgramConfig   `json:"-"`
}

func readConfig() Config {
	var configFile string
	var programName string
	var verbose nillaBool.NillaBool

	flag.StringVar(&configFile, "config", "config.json", "Path to config file")
	flag.StringVar(&programName, "program", "", "Name of program to launch")
	flag.Var(&verbose, "verbose", "Verbose output")

	flag.Parse()

	if programName == "" {
		fmt.Println("No program name provided")
		os.Exit(1)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("Failed to read config file:", err)
		os.Exit(1)
	}

	var c Config
	err = json.Unmarshal(data, &c)
	if err != nil {
		fmt.Println("Failed to parse config file:", err)
		os.Exit(1)
	}

	for _, program := range c.Programs {
		if strings.EqualFold(program.Name, programName) {
			c.Program = program
			break
		}
	}

	if c.Program.Name == "" {
		fmt.Println("Program not found in config file")
		os.Exit(1)
	}

	if !verbose.IsNull() {
		if v, isNil := verbose.Get(); !isNil {
			c.Verbose = v
		}
	}

	if c.Verbose {
		fmt.Println("Verbose!")
	}

	return c
}
