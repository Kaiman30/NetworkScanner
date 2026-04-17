package checker

import (
	"github.com/Kaiman30/NetworkChecker/pkg/windows"
)

// CheckMSIX проверяет режим MSI-X
func CheckMSIX(ctx *CheckContext) {
	pciPath := `SYSTEM\CurrentControlSet\Enum\PCI`

	pciDevices, err := windows.GetRegistrySubKeys(pciPath)
	if err != nil {
		ctx.Results.AddPassed("MSI-X Mode (не удалось проверить)")
		return
	}

	found := false
	for _, device := range pciDevices {
		msiPath := pciPath + `\` + device + `\Device Parameters\Interrupt Management\MessageSignaledInterruptProperties`

		if windows.RegistryKeyExists(msiPath) {
			val, err := windows.GetRegistryUint64(msiPath, "MSISupported")
			if err == nil && val == 1 {
				found = true
				break
			}
		}
	}

	if found {
		ctx.Results.AddFailed("MSI-X Mode (активно - оптимизация найдена)")
	} else {
		ctx.Results.AddPassed("MSI-X Mode")
	}
}
