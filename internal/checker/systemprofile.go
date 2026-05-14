package checker

import (
	"strconv"

	"github.com/Kaiman30/NetworkChecker/pkg/windows"
)

// defaultSystemProfileValues содержит приемлемые значения по умолчанию для параметров системного профиля
var defaultSystemProfileValues = map[string][]uint64{
	"NetworkThrottlingIndex": {0x0000000a, 0xffffff},
	"SystemResponsiveness":   {0x00000014},
}

func CheckSystemProfile(ctx *CheckContext) {
	sysProfilePath := `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Multimedia\SystemProfile`

	checks := []string{
		"NetworkThrottlingIndex",
		"SystemResponsiveness",
	}

	for _, checkName := range checks {
		val, err := windows.GetRegistryUint64(sysProfilePath, checkName)

		if err == nil {
			// Ключ существует - проверяем, является ли значение приемлемым по умолчанию
			if isDefaultSystemProfileValue(checkName, val) {
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

// isDefaultSystemProfileValue проверяет, является ли значение приемлемым по умолчанию
func isDefaultSystemProfileValue(paramName string, value uint64) bool {
	if defaults, exists := defaultSystemProfileValues[paramName]; exists {
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
