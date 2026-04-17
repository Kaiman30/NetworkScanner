package checker

import (
	"strconv"

	"github.com/Kaiman30/NetworkChecker/pkg/windows"
)

// CheckAFDParams проверяет параметры AFD (UDP)
func CheckAFDParams(ctx *CheckContext) {
	afdPath := `SYSTEM\CurrentControlSet\Services\AFD\Parameters`

	val, err := windows.GetRegistryUint64(afdPath, "FastSendDatagramThreshold")
	if err == nil {
		ctx.Results.AddFailed("FastSendDatagramThreshold = " + strconv.FormatUint(val, 10))
	} else {
		ctx.Results.AddPassed("FastSendDatagramThreshold")
	}
}
