package checker

import (
	"strconv"

	"github.com/Kaiman30/NetworkChecker/pkg/windows"
)

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
			// Ключ существует - значит есть настройка
			ctx.Results.AddFailed(checkName + " = " + strconv.FormatUint(val, 10))
		} else {
			// Ключа нет - настройки нет
			ctx.Results.AddPassed(checkName)
		}
	}
}
