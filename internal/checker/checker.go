package checker

import (
	"fmt"

	"github.com/Kaiman30/NetworkChecker/internal/models"
)

// RunAllChecks запускает все проверки
func RunAllChecks() (*models.Results, *FakerCheckResult) {
	results := models.NewResults()

	ctx := &CheckContext{
		Results: results,
	}

	fmt.Println("=== Network Settings Checker ===")
	fmt.Println("Поиск измененных сетевых настроек...")
	fmt.Println()

	CheckDriverKeys(ctx)
	CheckTCPIPParams(ctx)
	CheckAFDParams(ctx)
	CheckSystemProfile(ctx)
	CheckInterfaces(ctx)
	CheckNetshTCP(ctx)
	CheckMSIX(ctx)
	CheckQoS(ctx)

	fmt.Printf("\n=== РЕЗУЛЬТАТЫ ===\n")
	fmt.Printf("Настроек не обнаружено: %d\n", len(results.Passed))
	fmt.Printf("Обнаружены настройки: %d\n", len(results.Failed))

	if len(results.Failed) > 0 {
		fmt.Println("\nНайденные настройки:")
		for _, item := range results.Failed {
			fmt.Printf("  • %s\n", item)
		}
	}

	// Проверка на фейкеров
	fmt.Println("\n[Faker Check] Анализ сети на предмет фейкеров...")
	fakerResult := CheckFaker(ctx)

	return results, fakerResult
}

type CheckContext struct {
	Results *models.Results
}
