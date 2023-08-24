package winlaunch

import (
	"fmt"

	"github.com/emmaly/winlaunch/config"
	"github.com/emmaly/winlaunch/launch"
	"github.com/emmaly/winlaunch/window"
)

func Raise(program config.ProgramConfig, verbose bool) error {
	windowRaised := window.RaiseMatchedWindow(program.TitleMatch)
	if windowRaised {
		if verbose {
			fmt.Println("Window matched; raised.")
		}
		return nil
	}

	if verbose {
		fmt.Println("Window not matched; spawning new...")
	}
	launched := launch.LaunchProgram(program)

	if launched {
		if verbose {
			fmt.Println("Launched.")
		}
		return nil
	}

	if verbose {
		fmt.Println("Failed to launch.")
	}
	return fmt.Errorf("failed to launch")
}
