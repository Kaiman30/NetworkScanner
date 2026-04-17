package ui

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32         = windows.NewLazySystemDLL("user32.dll")
	procMessageBox = user32.NewProc("MessageBoxW")
)

const (
	MB_OK              = 0x00000000
	MB_ICONINFORMATION = 0x00000040
	MB_ICONWARNING     = 0x00000030
	MB_ICONERROR       = 0x00000010
)

// ShowMessageBox показывает окно сообщения
func showMessageBox(title, message string, icon uintptr) {
	procMessageBox.Call(0,
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(message))),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(title))),
		icon)
}

// ShowInfoMessage показывает информационное сообщение
func ShowInfoMessage(title, message string) {
	showMessageBox(title, message, MB_ICONINFORMATION)
}

// ShowErrorMessage показывает сообщение об ошибке
func ShowErrorMessage(title, message string) {
	showMessageBox(title, message, MB_ICONERROR)
}

// ShowWarningMessage показывает предупреждение
func ShowWarningMessage(title, message string) {
	showMessageBox(title, message, MB_ICONWARNING)
}
