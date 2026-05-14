package checker

import (
	"strconv"

	"github.com/Kaiman30/NetworkChecker/pkg/windows"
)

// defaultTCPIPValues содержит приемлемые значения по умолчанию для параметров TCP/IP
// Если параметр имеет одно из этих значений, он не считается измененной настройкой
var defaultTCPIPValues = map[string][]uint64{
	"Tcp1323Opts": {0, 1, 2},
	"DefaultTTL":  {64, 128},
	"SackOpts":    {1},
}

// CheckTCPIPParams проверяет параметры TCP/IP
func CheckTCPIPParams(ctx *CheckContext) {
	tcpipPath := `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`

	checks := []string{
		"SackOpts",
		"DisableTaskOffload",
		"EnableWsd",
		"Tcp1323Opts",
		"DefaultTTL",
		"EnablePMTUDiscovery",
		"EnablePMTUBHDetect",
		"GlobalMaxTcpWindowSize",
		"TcpMaxDataRetransmissions",
	}

	for _, checkName := range checks {
		val, err := windows.GetRegistryUint64(tcpipPath, checkName)

		if err == nil {
			// Ключ существует - проверяем, является ли значение приемлемым по умолчанию
			if isDefaultTcpIPValue(checkName, val) {
				ctx.Results.AddPassed(checkName)
			} else {
				// Значение измененное - добавляем в failed
				ctx.Results.AddFailed(checkName + " = " + strconv.FormatUint(val, 10))
			}
		} else {
			// Ключа нет - настройки нет
			ctx.Results.AddPassed(checkName)
		}
	}
}

// isDefaultTcpIPValue проверяет, является ли значение приемлемым по умолчанию
func isDefaultTcpIPValue(paramName string, value uint64) bool {
	if defaults, exists := defaultTCPIPValues[paramName]; exists {
		for _, defaultVal := range defaults {
			if value == defaultVal {
				return true
			}
		}
		return false
	}
	// Если параметр не в списке дефолтных значений, то любое значение считается измененным
	return false
}
