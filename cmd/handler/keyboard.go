//go:build windows

// go build -ldflags -H=windowsgui

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"github.com/emmaly/winlaunch"
	"github.com/emmaly/winlaunch/config"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"github.com/moutend/go-hook/pkg/win32"
)

func main() {
	cfg := config.ReadConfig()

	println("RUNNING")
	println(run(cfg, cfg.Verbose))
}

func run(cfg config.Config, verbose bool) error {
	keys := make(map[string][]config.ProgramConfig)
	for _, p := range cfg.Programs {
		if p.Hotkey == "" {
			continue
		}
		if _, ok := keys[p.Hotkey]; !ok {
			keys[p.Hotkey] = []config.ProgramConfig{p}
		} else {
			keys[p.Hotkey] = append(keys[p.Hotkey], p)
		}
	}

	hookHandler := func(c chan<- types.KeyboardEvent) types.HOOKPROC {
		return func(code int32, wParam, lParam uintptr) uintptr {
			if lParam != 0 {
				message := types.Message(wParam)
				kbdStruct := *(*types.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

				if message == types.WM_KEYDOWN || message == types.WM_SYSKEYDOWN || message == types.WM_KEYUP || message == types.WM_SYSKEYUP {
					if verbose {
						fmt.Printf("Checking %s [%s]\n", message.String(), kbdStruct.VKCode.String())
					}
					if _, ok := keys[kbdStruct.VKCode.String()]; ok {
						c <- types.KeyboardEvent{
							Message:         message,
							KBDLLHOOKSTRUCT: kbdStruct,
						}
						return 1
					}
				}
			}

			return win32.CallNextHookEx(0, code, wParam, lParam)
		}
	}

	keyboardChan := make(chan types.KeyboardEvent, 100)
	if verbose {
		println("KEYBOARD INSTALLING")
	}
	if err := keyboard.Install(hookHandler, keyboardChan); err != nil {
		if verbose {
			println("KEYBOARD INSTALL ERROR")
		}
		return err
	}
	defer keyboard.Uninstall()
	if verbose {
		println("KEYBOARD INSTALLED")
	}

	signalChan := make(chan os.Signal, 1)
	if verbose {
		println("SIGNAL INSTALLING")
	}
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	if verbose {
		println("SIGNAL INSTALLED")
	}

	for {
		select {
		case <-time.After(1 * time.Minute):
			if verbose {
				println("1-minute timeout")
			}
			return nil
		case <-signalChan:
			if verbose {
				println("signal received")
			}
			return nil
		case event := <-keyboardChan:
			if event.Message == types.WM_KEYDOWN || event.Message == types.WM_SYSKEYDOWN {
				fmt.Printf("Event: %s [%s]\n", event.Message.String(), event.VKCode.String())
				if _, ok := keys[event.VKCode.String()]; ok {
					for _, program := range keys[event.VKCode.String()] {
						fmt.Printf("Raising/Launching: %s\n", program.Name)
						err := winlaunch.Raise(program, cfg.Verbose)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
	}
}
