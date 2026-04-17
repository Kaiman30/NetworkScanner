package ui

import (
	"fmt"
	"os"
)

// UpdateProgress обновляет индикатор прогресса (через заголовок окна)
func UpdateProgress(stepName string, current, total int) {
	percent := int(float64(current) / float64(total) * 100)
	title := fmt.Sprintf("Network Checker - %d%% - %s", percent, stepName)
	os.Setenv("GUI_TITLE", title)
}
