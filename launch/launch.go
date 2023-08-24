package launch

import (
	"os/exec"

	"github.com/emmaly/winlaunch/config"
)

func LaunchProgram(program config.ProgramConfig) bool {
	cmd := exec.Command(program.Command, program.Args...)
	if program.StartIn != "" {
		cmd.Dir = program.StartIn
	}

	err := cmd.Start()
	return err == nil
}
