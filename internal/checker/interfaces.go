package checker

import (
	"fmt"
	"strconv"

	"github.com/Kaiman30/NetworkChecker/pkg/windows"
)

// CheckInterfaces проверяет настройки интерфейсов (Nagle/ACK)
func CheckInterfaces(ctx *CheckContext) {
	interfacesPath := `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces`

	interfaceKeys, err := windows.GetRegistrySubKeys(interfacesPath)
	if err != nil {
		ctx.Results.AddFailed("Не удалось прочитать интерфейсы")
		return
	}

	checks := map[string]string{
		"TCPNoDelay":      "Nagle",
		"TcpAckFrequency": "ACK Frequency",
		"TcpDelAckTicks":  "ACK Delay",
	}

	for checkName, checkDesc := range checks {
		found := false
		var foundValue uint64

		for _, subKey := range interfaceKeys {
			subKeyPath := interfacesPath + "\\" + subKey

			// Пропускаем интерфейсы без IP
			dhcpIP, _ := windows.GetRegistryString(subKeyPath, "DhcpIPAddress")
			ipAddress, _ := windows.GetRegistryString(subKeyPath, "IPAddress")

			if dhcpIP == "" && (ipAddress == "" || ipAddress == "0.0.0.0") {
				continue
			}

			val, err := windows.GetRegistryUint64(subKeyPath, checkName)
			if err == nil {
				found = true
				foundValue = val
				break
			}
		}

		if found {
			ctx.Results.AddFailed(fmt.Sprintf("%s (%s) = %s", checkName, checkDesc, strconv.FormatUint(foundValue, 10)))
		} else {
			ctx.Results.AddPassed(checkName)
		}
	}
}
