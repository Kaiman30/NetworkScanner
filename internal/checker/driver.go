package checker

import (
	"github.com/Kaiman30/NetworkChecker/pkg/windows"
)

// CheckDriverKeys проверяет ключи драйверов сетевых адаптеров
func CheckDriverKeys(ctx *CheckContext) {
	adapterClassPath := `SYSTEM\CurrentControlSet\Control\Class\{4d36e972-e325-11ce-bfc1-08002be10318}`

	adapterKeys, err := windows.GetRegistrySubKeys(adapterClassPath)
	if err != nil {
		ctx.Results.AddFailed("Не удалось прочитать ключи драйверов")
		return
	}

	checks := map[string]uint64{
		"*NdisDeviceType": 1,
		"*JumboPacket":    1514,
	}

	for checkName, expected := range checks {
		found := false
		for _, subKey := range adapterKeys {
			if len(subKey) == 4 && subKey[0] >= '0' && subKey[0] <= '9' {
				val, err := windows.GetRegistryUint64(adapterClassPath+"\\"+subKey, checkName)
				if err == nil && val == expected {
					found = true
					break
				}
			}
		}
		if found {
			ctx.Results.AddFailed(checkName + " (найдена оптимизация драйвера)")
		} else {
			ctx.Results.AddPassed(checkName)
		}
	}
}
