package windows

import (
	"golang.org/x/sys/windows"
)

var (
	user32         = windows.NewLazySystemDLL("user32.dll")
	kernel32       = windows.NewLazySystemDLL("kernel32.dll")
	procGetConsole = kernel32.NewProc("GetConsoleWindow")
	procShowWindow = user32.NewProc("ShowWindow")
)

const (
	SW_HIDE = 0
)

// HideConsole скрывает консольное окно
func HideConsole() {
	consoleWindow, _, _ := procGetConsole.Call()
	if consoleWindow != 0 {
		procShowWindow.Call(consoleWindow, SW_HIDE)
	}
}

// GetConsoleHandle возвращает хендл консоли
func GetConsoleHandle() uintptr {
	handle, _, _ := procGetConsole.Call()
	return handle
}
