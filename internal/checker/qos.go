package checker

import (
	"strconv"

	"github.com/Kaiman30/NetworkChecker/pkg/windows"
)

// CheckQoS проверяет настройки QoS Bandwidth
func CheckQoS(ctx *CheckContext) {
	qosPath := `SOFTWARE\Policies\Microsoft\Windows\Psched`

	val, err := windows.GetRegistryUint64(qosPath, "NonBestEffortLimit")
	if err == nil {
		ctx.Results.AddFailed("NonBestEffortLimit = " + strconv.FormatUint(val, 10))
	} else {
		ctx.Results.AddPassed("NonBestEffortLimit")
	}
}
