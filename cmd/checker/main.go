package main

import (
	"github.com/Kaiman30/NetworkChecker/internal/checker"
	"github.com/Kaiman30/NetworkChecker/internal/ui"
	"github.com/Kaiman30/NetworkChecker/internal/web"
	"github.com/Kaiman30/NetworkChecker/pkg/windows"
)

func main() {
	// Проверка прав администратора
	if !windows.IsAdmin() {
		ui.ShowErrorMessage("Administrator Required", "Этот скрипт должен запускаться от имени администратора!")
		return
	}

	// Скрываем окно терминала
	windows.HideConsole()

	// Стартовое сообщение
	ui.ShowInfoMessage("Network Settings Checker", "Начало анализа...\nОтчет откроется в браузере")

	results := checker.RunAllChecks()

	// Запус сервера с результатами
	web.StartServer(results)
}
