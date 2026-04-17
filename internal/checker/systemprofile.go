package checker

import (
	"strconv"

	"github.com/Kaiman30/NetworkChecker/pkg/windows"
)

// CheckSystemProfile проверяет системный профиль (троттлинг)
func CheckSystemProfile(ctx *CheckContext) {
	sysProfilePath := `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Multimedia\SystemProfile`

	// NetworkThrottlingIndex
	val, err := windows.GetRegistryUint64(sysProfilePath, "NetworkThrottlingIndex")
	if err == nil {
		ctx.Results.AddFailed("NetworkThrottlingIndex = " + formatThrottleValue(val))
	} else {
		ctx.Results.AddPassed("NetworkThrottlingIndex")
	}

	// SystemResponsiveness
	val, err = windows.GetRegistryUint64(sysProfilePath, "SystemResponsiveness")
	if err == nil {
		ctx.Results.AddFailed("SystemResponsiveness = " + strconv.FormatUint(val, 10))
	} else {
		ctx.Results.AddPassed("SystemResponsiveness")
	}
}

func formatThrottleValue(val uint64) string {
	if val == 0xFFFFFFFF {
		return "отключен"
	}
	return strconv.FormatUint(val, 10)
}
