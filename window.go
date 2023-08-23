package main

import (
	"regexp"
	"unsafe"
)

/*
#include <windows.h>
#include <stdbool.h>

// Forward declarations of functions to avoid multiple definitions
BOOL CheckWindowTitle(HWND hwnd, char* outStr, int maxLen);
BOOL SetAsForegroundWindow(HWND hwnd);
BOOL CALLBACK EnumWindowsProc(HWND hwnd, LPARAM lParam);
*/
import "C"

//export CheckWindowTitle
func CheckWindowTitle(hwnd C.HWND, outStr *C.char, maxLen C.int) C.BOOL {
	return C.GetWindowTextA(hwnd, outStr, maxLen)
}

//export SetAsForegroundWindow
func SetAsForegroundWindow(hwnd C.HWND) C.BOOL {
	if (C.IsIconic(hwnd)) != 0 {
		C.ShowWindow(hwnd, C.SW_RESTORE)
	}
	return C.SetForegroundWindow(hwnd)
}

//export EnumWindowsProc
func EnumWindowsProc(hwnd C.HWND, lParam C.LPARAM) C.int {
	// return 1 if no match is found, 0 if a match is found

	//nolint:unsafeptr // ptr is safe
	titleMatch := C.GoString((*C.char)(unsafe.Pointer(uintptr(lParam))))

	const maxLen = 256
	var buffer [maxLen]C.char

	if CheckWindowTitle(hwnd, &buffer[0], maxLen) == 0 {
		return 1 // no match
	}

	title := C.GoString(&buffer[0])
	matched, err := regexp.MatchString(titleMatch, title)
	if err != nil {
		panic(err)
	}

	if matched {
		if verboseOutput {
			println("Matched window: ", title)
		}
		SetAsForegroundWindow(hwnd)
		return 0 // match
	}

	return 1 // no match
}

func raiseMatchedWindow(titleMatch string) bool {
	titleMatchCString := C.CString(titleMatch)
	defer C.free(unsafe.Pointer(titleMatchCString))
	result := C.EnumWindows(C.WNDENUMPROC(C.EnumWindowsProc), C.LPARAM(uintptr(unsafe.Pointer(titleMatchCString))))
	return result == 0
}
