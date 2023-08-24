//go:build !windows

package window

func RaiseMatchedWindow(titleMatch string) bool {
	panic("not implemented on non-Windows platforms")
	return false
}
