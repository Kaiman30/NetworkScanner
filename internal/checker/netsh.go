package checker

import (
	"os/exec"
	"strings"
)

// CheckNetshTCP проверяет настройки netsh tcp global
func CheckNetshTCP(ctx *CheckContext) {
	cmd := exec.Command("netsh", "int", "tcp", "show", "global")
	output, err := cmd.Output()
	if err != nil {
		ctx.Results.AddFailed("Не удалось выполнить netsh")
		return
	}

	outStr := strings.ToLower(string(output))
	lines := strings.Split(outStr, "\n")

	// Проверяем каждую настройку
	settings := map[string]string{
		"receive window":  "Receive Window Auto-Tuning",
		"receive segment": "Receive Segment Coalescing",
		"ecn":             "ECN Capability",
	}

	for keyword, settingName := range settings {
		for _, line := range lines {
			if strings.Contains(line, keyword) {
				if strings.Contains(line, "enabled") {
					ctx.Results.AddFailed(settingName + " = enabled")
				} else if strings.Contains(line, "disabled") {
					ctx.Results.AddPassed(settingName)
				}
				break
			}
		}
	}
}
