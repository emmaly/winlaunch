package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/emmaly/winlaunch/nillaBool"
)

type ProgramConfig struct {
	Name       string   `json:"name"`
	TitleMatch string   `json:"titleMatch"`
	Command    string   `json:"command"`
	Args       []string `json:"args"`
	StartIn    string   `json:"startIn"`
	Hotkey     string   `json:"hotkey"`
}

type Config struct {
	Verbose  bool            `json:"verbose"`
	Programs []ProgramConfig `json:"programs"`
	Args     []string        `json:"-"`
}

func ReadConfig() Config {
	var configFile string
	var verbose nillaBool.NillaBool

	flag.StringVar(&configFile, "config", "config.json", "Path to config file")
	flag.Var(&verbose, "verbose", "Verbose output")
	flag.Parse()

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

	c.Args = flag.Args()

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
