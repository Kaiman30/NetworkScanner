package web

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Kaiman30/NetworkChecker/internal/models"
	"github.com/Kaiman30/NetworkChecker/internal/ui"
)

// StartServer создает HTML файл и открывает его в браузере
func StartServer(results *models.Results) {
	// Создаем временный HTML файл
	tempDir := os.TempDir()
	htmlFile := filepath.Join(tempDir, "network_checker_report.html")

	// Генерируем JSON данные
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		ui.ShowErrorMessage("Ошибка", "Не удалось сгенерировать отчет: "+err.Error())
		return
	}

	// Для отладки - выводим данные в консоль
	fmt.Printf("Results JSON: %s\n", string(data))
	fmt.Printf("Passed count: %d\n", len(results.Passed))
	fmt.Printf("Failed count: %d\n", len(results.Failed))

	// Создаем HTML контент
	htmlContent := createHTMLContent(string(data))

	// Записываем файл
	err = os.WriteFile(htmlFile, []byte(htmlContent), 0644)
	if err != nil {
		ui.ShowErrorMessage("Ошибка", "Не удалось сохранить отчет: "+err.Error())
		return
	}

	// Открываем файл в браузере
	openFileInBrowser(htmlFile)
}

func createHTMLContent(jsonData string) string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Network Settings Checker - Report</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
            min-height: 100vh;
            padding: 20px;
            color: #eee;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        h1 {
            text-align: center;
            color: #c792ea;
            margin-bottom: 30px;
            font-size: 2.5em;
            text-shadow: 0 0 10px rgba(199,146,234,0.5);
        }
        .timestamp {
            text-align: center;
            color: #888;
            margin-bottom: 30px;
            font-size: 0.9em;
        }
        .stats {
            display: flex;
            gap: 20px;
            justify-content: center;
            margin-bottom: 40px;
            flex-wrap: wrap;
        }
        .stat-card {
            background: rgba(255,255,255,0.1);
            border-radius: 15px;
            padding: 20px 40px;
            text-align: center;
            backdrop-filter: blur(10px);
            min-width: 200px;
            transition: transform 0.3s;
        }
        .stat-card:hover {
            transform: translateY(-5px);
        }
        .stat-number {
            font-size: 3em;
            font-weight: bold;
        }
        .stat-label {
            font-size: 0.9em;
            opacity: 0.8;
            margin-top: 5px;
        }
        .good .stat-number { color: #4caf50; }
        .bad .stat-number { color: #f44336; }
        .section {
            background: rgba(255,255,255,0.05);
            border-radius: 15px;
            padding: 20px;
            margin-bottom: 20px;
            transition: all 0.3s;
        }
        .section:hover {
            background: rgba(255,255,255,0.08);
        }
        .section h2 {
            margin-bottom: 15px;
            padding-bottom: 10px;
            border-bottom: 2px solid;
        }
        .section.good h2 { border-color: #4caf50; color: #4caf50; }
        .section.bad h2 { border-color: #f44336; color: #f44336; }
        .item-list {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
        }
        .item {
            background: rgba(255,255,255,0.1);
            padding: 8px 15px;
            border-radius: 20px;
            font-family: monospace;
            font-size: 0.9em;
            transition: all 0.3s;
        }
        .item:hover {
            background: rgba(255,255,255,0.2);
            transform: scale(1.05);
        }
        .empty {
            color: #666;
            font-style: italic;
            padding: 20px;
            text-align: center;
        }
        footer {
            text-align: center;
            margin-top: 30px;
            opacity: 0.6;
            font-size: 0.8em;
        }
        .note {
            background: rgba(255,255,255,0.05);
            border-left: 4px solid #c792ea;
            padding: 15px;
            margin-top: 20px;
            border-radius: 5px;
            font-size: 0.9em;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔍 Network Settings Checker</h1>
        <div class="timestamp">
            Report generated: <span id="timestamp"></span>
        </div>
        <div class="stats">
            <div class="stat-card good">
                <div class="stat-number" id="passedCount">0</div>
                <div class="stat-label">✅ Настроек не выявлено</div>
            </div>
            <div class="stat-card bad">
                <div class="stat-number" id="failedCount">0</div>
                <div class="stat-label">❌ Обнаружены настройки</div>
            </div>
        </div>
        <div class="section good">
            <h2>✅ Настройки не обнаружены (оптимально)</h2>
            <div id="passedList" class="item-list"></div>
        </div>
        <div class="section bad">
            <h2>❌ Обнаруженные настройки (требуют внимания)</h2>
            <div id="failedList" class="item-list"></div>
        </div>
        <div class="note">
            💡 <strong>Пояснение:</strong><br>
            • "Настроек не выявлено" - параметры в стандартном состоянии<br>
            • "Обнаружены настройки" - параметры изменены от стандартных
        </div>
        <footer>
            Network Settings Checker by Kaiman4ik :3
        </footer>
    </div>
    <script>
        const results = ` + jsonData + `;
        
        results.passed = results.passed || [];
        results.failed = results.failed || [];
        
        document.getElementById('passedCount').textContent = results.passed.length;
        document.getElementById('failedCount').textContent = results.failed.length;
        
        function renderList(containerId, items, emptyMessage) {
            const container = document.getElementById(containerId);
            if (!items || items.length === 0) {
                container.innerHTML = '<div class="empty">' + emptyMessage + '</div>';
                return;
            }
            container.innerHTML = items.map(item => '<div class="item">' + escapeHtml(item) + '</div>').join('');
        }
        
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
        
        renderList('passedList', results.passed, '✨ Все параметры в стандартном состоянии');
        renderList('failedList', results.failed, '🔧 Нет измененных параметров');
        
        document.getElementById('timestamp').textContent = new Date().toLocaleString();
        
        console.log('Passed:', results.passed);
        console.log('Failed:', results.failed);
    </script>
</body>
</html>`
}

// openFileInBrowser открывает файл в браузере
func openFileInBrowser(filePath string) {
	var cmd *exec.Cmd

	// Конвертируем путь в file:// URL
	fileURL := "file:///" + filepath.ToSlash(filePath)

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", fileURL)
		if err := cmd.Run(); err != nil {
			// Пробуем альтернативный метод
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", fileURL)
			cmd.Run()
		}
	default:
		cmd = exec.Command("xdg-open", fileURL)
		cmd.Run()
	}
}
