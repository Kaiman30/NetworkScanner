package web

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Kaiman30/NetworkChecker/internal/checker"
	"github.com/Kaiman30/NetworkChecker/internal/models"
)

// StartServer создает HTML файл и открывает его в браузере
func StartServer(results *models.Results, fakerResult *checker.FakerCheckResult) {
	// Создаем временный HTML файл
	tempDir := os.TempDir()
	htmlFile := filepath.Join(tempDir, "network_checker_report.html")

	// Генерируем JSON данные
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("Error: Не удалось сгенерировать отчет: %v\n", err)
		return
	}

	// Генерируем JSON для фейкер данных
	fakerData, err := json.MarshalIndent(fakerResult, "", "  ")
	if err != nil {
		fmt.Printf("Error: Не удалось сгенерировать отчет фейкеров: %v\n", err)
		return
	}

	// Для отладки - выводим данные в консоль
	fmt.Printf("Results JSON: %s\n", string(data))
	fmt.Printf("Passed count: %d\n", len(results.Passed))
	fmt.Printf("Failed count: %d\n", len(results.Failed))

	// Создаем HTML контент с табами
	htmlContent := createHTMLContent(string(data), string(fakerData))

	// Записываем файл
	err = os.WriteFile(htmlFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error: Не удалось сохранить отчет: %v\n", err)
		return
	}

	// Открываем файл в браузере
	openFileInBrowser(htmlFile)
}

func createHTMLContent(jsonData string, fakerData string) string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Network Settings Checker - Report</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }

        :root {
            --primary: #6366f1;
            --primary-dark: #4f46e5;
            --primary-light: #818cf8;
            --success: #10b981;
            --danger: #ef4444;
            --warning: #f59e0b;
            --dark-bg: #0f172a;
            --card-bg: #1e293b;
            --border-color: #334155;
            --text-primary: #f1f5f9;
            --text-secondary: #cbd5e1;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', sans-serif;
            background: linear-gradient(135deg, var(--dark-bg) 0%, #1a2332 100%);
            min-height: 100vh;
            color: var(--text-primary);
            line-height: 1.6;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
            padding: 40px 20px;
        }

        .header {
            text-align: center;
            margin-bottom: 50px;
        }

        .header h1 {
            font-size: 3em;
            font-weight: 800;
            background: linear-gradient(135deg, var(--primary-light) 0%, #06b6d4 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            margin-bottom: 10px;
            letter-spacing: -1px;
        }

        .header p {
            color: var(--text-secondary);
            font-size: 1.1em;
        }

        .tabs {
            display: flex;
            gap: 15px;
            margin-bottom: 40px;
            border-bottom: 2px solid var(--border-color);
            padding-bottom: 0;
        }

        .tab-button {
            background: transparent;
            border: none;
            color: var(--text-secondary);
            padding: 15px 25px;
            cursor: pointer;
            font-size: 1.05em;
            font-weight: 600;
            transition: all 0.3s ease;
            border-bottom: 3px solid transparent;
            position: relative;
            bottom: -2px;
        }

        .tab-button:hover {
            color: var(--text-primary);
        }

        .tab-button.active {
            color: var(--primary-light);
            border-bottom-color: var(--primary-light);
        }

        .tab-content {
            display: none;
            animation: fadeIn 0.3s ease-in;
        }

        .tab-content.active {
            display: block;
        }

        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }

        .timestamp {
            text-align: center;
            color: var(--text-secondary);
            margin-bottom: 30px;
            font-size: 0.95em;
        }

        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-bottom: 40px;
        }

        .stat-card {
            background: linear-gradient(135deg, var(--card-bg) 0%, rgba(30, 41, 59, 0.5) 100%);
            border: 1px solid var(--border-color);
            border-radius: 16px;
            padding: 30px 25px;
            text-align: center;
            transition: all 0.3s ease;
            cursor: pointer;
        }

        .stat-card:hover {
            border-color: var(--primary);
            transform: translateY(-5px);
            box-shadow: 0 20px 25px -5px rgba(99, 102, 241, 0.1);
        }

        .stat-number {
            font-size: 2.5em;
            font-weight: 900;
            margin-bottom: 8px;
        }

        .stat-label {
            font-size: 0.95em;
            color: var(--text-secondary);
            font-weight: 500;
        }

        .stat-card.good .stat-number { color: var(--success); }
        .stat-card.bad .stat-number { color: var(--danger); }

        .section {
            background: linear-gradient(135deg, var(--card-bg) 0%, rgba(30, 41, 59, 0.5) 100%);
            border: 1px solid var(--border-color);
            border-radius: 16px;
            padding: 30px;
            margin-bottom: 25px;
            transition: all 0.3s ease;
        }

        .section:hover {
            border-color: var(--primary);
        }

        .section h2 {
            margin-bottom: 20px;
            padding-bottom: 15px;
            font-size: 1.4em;
            font-weight: 700;
            border-bottom: 2px solid;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .section.good h2 {
            border-color: var(--success);
            color: var(--success);
        }
        .section.bad h2 {
            border-color: var(--danger);
            color: var(--danger);
        }

        .item-list {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
            gap: 12px;
        }

        .item {
            background: rgba(148, 163, 184, 0.1);
            border: 1px solid var(--border-color);
            padding: 12px 16px;
            border-radius: 10px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 0.9em;
            transition: all 0.3s ease;
            word-break: break-word;
        }

        .item:hover {
            background: rgba(99, 102, 241, 0.15);
            border-color: var(--primary);
            transform: translateX(4px);
        }

        .empty {
            color: var(--text-secondary);
            font-style: italic;
            padding: 30px 20px;
            text-align: center;
            grid-column: 1 / -1;
        }

        footer {
            text-align: center;
            margin-top: 60px;
            padding-top: 30px;
            border-top: 1px solid var(--border-color);
            color: var(--text-secondary);
            font-size: 0.9em;
        }

        footer a {
            color: var(--primary-light);
            text-decoration: none;
            font-weight: 600;
            transition: color 0.3s;
        }

        footer a:hover {
            color: var(--primary-light);
            text-decoration: underline;
        }

        .note {
            background: linear-gradient(135deg, rgba(99, 102, 241, 0.15) 0%, rgba(6, 182, 212, 0.1) 100%);
            border-left: 4px solid var(--primary-light);
            border-radius: 10px;
            padding: 20px;
            margin-top: 25px;
            font-size: 0.95em;
        }

        .note strong {
            color: var(--primary-light);
        }

        .note ul {
            margin: 12px 0 0 20px;
        }

        .note li {
            margin: 6px 0;
            color: var(--text-secondary);
        }

        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            font-size: 0.95em;
        }

        th {
            background: rgba(99, 102, 241, 0.2);
            padding: 15px;
            text-align: left;
            font-weight: 700;
            border-bottom: 2px solid var(--primary);
            color: var(--primary-light);
        }

        td {
            padding: 12px 15px;
            border-bottom: 1px solid var(--border-color);
        }

        tr:hover {
            background: rgba(99, 102, 241, 0.1);
        }

        .warning-box {
            background: linear-gradient(135deg, rgba(239, 68, 68, 0.15) 0%, rgba(239, 68, 68, 0.05) 100%);
            border-left: 4px solid var(--danger);
            border-radius: 10px;
            padding: 20px;
            margin: 20px 0;
            color: #fca5a5;
            font-weight: 500;
        }

        .success-box {
            background: linear-gradient(135deg, rgba(16, 185, 129, 0.15) 0%, rgba(16, 185, 129, 0.05) 100%);
            border-left: 4px solid var(--success);
            border-radius: 10px;
            padding: 20px;
            margin: 20px 0;
            color: #86efac;
            font-weight: 500;
        }

        h3 {
            margin-top: 25px;
            margin-bottom: 15px;
            font-size: 1.15em;
            color: var(--primary-light);
        }

        ul {
            margin-left: 20px;
            color: var(--text-secondary);
        }

        li {
            margin-bottom: 8px;
            line-height: 1.8;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 Network Settings Checker</h1>
            <p>Advanced Network Diagnostics Report</p>
        </div>

        <div class="tabs">
            <button class="tab-button active" onclick="switchTab('network')">🖥️ Network Settings</button>
            <button class="tab-button" onclick="switchTab('faker')">🚨 Faker Detection</button>
        </div>

        <!-- Network Settings Tab -->
        <div id="network" class="tab-content active">
            <div class="timestamp">
                Report generated: <span id="timestamp"></span>
            </div>
            <div class="stats">
                <div class="stat-card good">
                    <div class="stat-number" id="passedCount">0</div>
                    <div class="stat-label">✅ Settings OK</div>
                </div>
                <div class="stat-card bad">
                    <div class="stat-number" id="failedCount">0</div>
                    <div class="stat-label">⚠️ Issues Found</div>
                </div>
            </div>
            <div class="section good">
                <h2>✅ Optimal Configuration</h2>
                <div id="passedList" class="item-list"></div>
            </div>
            <div class="section bad">
                <h2>⚠️ Modified Settings</h2>
                <div id="failedList" class="item-list"></div>
            </div>
            <div class="note">
                <strong>💡 Information:</strong>
                <ul>
                    <li><strong>Optimal Configuration</strong> - Parameters in default state</li>
                    <li><strong>Modified Settings</strong> - Parameters changed from defaults</li>
                </ul>
            </div>
        </div>

        <!-- Faker Detection Tab -->
        <div id="faker" class="tab-content">
            <div id="fakerContent"></div>
        </div>

        <footer id="footer">
            <strong>Network Settings Checker</strong> by <a href="https://github.com/Kaiman30" style="color: var(--primary-light); text-decoration: none;">Kaiman4ik</a> :3
        </footer>
    </div>

    <script>
        const results = ` + jsonData + `;
        const fakerResults = ` + fakerData + `;

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

        renderList('passedList', results.passed, '✨ Default settings detected - no issues found!');
        renderList('failedList', results.failed, '🔧 No modified parameters found');

        document.getElementById('timestamp').textContent = new Date().toLocaleString();

        function switchTab(tabName) {
            // Hide all tabs
            document.querySelectorAll('.tab-content').forEach(tab => {
                tab.classList.remove('active');
            });
            document.querySelectorAll('.tab-button').forEach(btn => {
                btn.classList.remove('active');
            });

            // Show selected tab
            document.getElementById(tabName).classList.add('active');
            event.target.classList.add('active');

            // Update footer
            const footer = document.getElementById('footer');
            if (tabName === 'faker') {
                footer.innerHTML = 'Network Settings Checker by Kaiman4ik :3 | Faker Detection by <a href="https://github.com/praiselily" style="color: #c792ea; text-decoration: none;">@praiselily</a>';
            } else {
                footer.innerHTML = 'Network Settings Checker by Kaiman4ik :3';
            }
        }

        // Render Faker tab content
        function renderFakerTab() {
            const container = document.getElementById('fakerContent');
            let html = '<div class="timestamp">Report generated: ' + new Date().toLocaleString() + '</div>';

            if (fakerResults.fakerDetected) {
                html += '<div class="warning-box"><strong>⚠️ FAKER INDICATORS DETECTED!</strong></div>';
                html += '<h2 style="color: #f44336;">🚨 Faker Indicators:</h2>';
                if (fakerResults.fakerIndicators && fakerResults.fakerIndicators.length > 0) {
                    html += '<ul>';
                    fakerResults.fakerIndicators.forEach(indicator => {
                        html += '<li>' + escapeHtml(indicator) + '</li>';
                    });
                    html += '</ul>';
                }
            } else {
                html += '<div class="success-box"><strong>✅ No faker indicators detected</strong></div>';
            }

            // Suspicious Activities
            html += '<h2>🔍 Suspicious Activities</h2>';
            if (fakerResults.suspiciousActivities && fakerResults.suspiciousActivities.length > 0) {
                html += '<ul>';
                fakerResults.suspiciousActivities.forEach(activity => {
                    html += '<li>' + escapeHtml(activity) + '</li>';
                });
                html += '</ul>';
            } else {
                html += '<p class="empty">No suspicious activities detected</p>';
            }

            // Current Connection
            if (fakerResults.currentConnection) {
                html += '<h2>📡 Current Connection</h2>';
                html += '<table><tr><th>Property</th><th>Value</th></tr>';
                html += '<tr><td>SSID</td><td>' + escapeHtml(fakerResults.currentConnection.SSID) + '</td></tr>';
                html += '<tr><td>State</td><td>' + escapeHtml(fakerResults.currentConnection.State) + '</td></tr>';
                html += '<tr><td>BSSID</td><td>' + escapeHtml(fakerResults.currentConnection.BSSID) + '</td></tr>';
                html += '<tr><td>Network Type</td><td>' + escapeHtml(fakerResults.currentConnection.NetworkType) + '</td></tr>';
                html += '<tr><td>Channel</td><td>' + escapeHtml(fakerResults.currentConnection.Channel) + '</td></tr>';
                html += '<tr><td>Signal</td><td>' + escapeHtml(fakerResults.currentConnection.Signal) + '</td></tr>';
                html += '</table>';

                if (fakerResults.currentConnection.hotspotIndicators && fakerResults.currentConnection.hotspotIndicators.length > 0) {
                    html += '<h3>Hotspot Indicators:</h3><ul>';
                    fakerResults.currentConnection.hotspotIndicators.forEach(indicator => {
                        html += '<li>' + escapeHtml(indicator) + '</li>';
                    });
                    html += '</ul>';
                }
            }

            // Mobile Hotspot Service
            html += '<h2>📱 Mobile Hotspot Service</h2>';
            if (fakerResults.mobileHotspotActive) {
                html += '<div class="warning-box"><strong>RUNNING</strong> - This is a faker indicator!</div>';
            } else {
                html += '<div class="success-box"><strong>STOPPED</strong></div>';
            }

            // Hosted Network
            html += '<h2>🌐 Hosted Network Status</h2>';
            if (fakerResults.hostedNetworkActive) {
                html += '<div class="warning-box">ACTIVE - SSID: ' + escapeHtml(fakerResults.hostedNetworkSSID) + ', Clients: ' + fakerResults.hostedNetworkClients + '</div>';
            } else {
                html += '<div class="success-box">INACTIVE</div>';
            }

            // Hotspot Profiles
            if (fakerResults.hotspotProfiles && fakerResults.hotspotProfiles.length > 0) {
                html += '<h2>📶 Detected Hotspot Profiles</h2>';
                html += '<ul>';
                fakerResults.hotspotProfiles.forEach(profile => {
                    html += '<li>' + escapeHtml(profile) + '</li>';
                });
                html += '</ul>';
            }

            // Virtual Adapters
            if (fakerResults.virtualAdapters && fakerResults.virtualAdapters.length > 0) {
                html += '<h2>💻 Virtual Network Adapters</h2>';
                html += '<table><tr><th>Description</th></tr>';
                fakerResults.virtualAdapters.forEach(adapter => {
                    html += '<tr><td>' + escapeHtml(adapter.description) + '</td></tr>';
                });
                html += '</table>';
            }

            // Connected Devices
            if (fakerResults.connectedDevices && fakerResults.connectedDevices.length > 0) {
                html += '<h2>🖥️ Connected Devices</h2>';
                html += '<table><tr><th>IP Address</th><th>MAC Address</th><th>Type</th></tr>';
                fakerResults.connectedDevices.forEach(device => {
                    html += '<tr><td>' + escapeHtml(device.ip) + '</td><td>' + escapeHtml(device.mac) + '</td><td>' + escapeHtml(device.type) + '</td></tr>';
                });
                html += '</table>';
            }

            container.innerHTML = html;
        }

        renderFakerTab();

        console.log('Passed:', results.passed);
        console.log('Failed:', results.failed);
        console.log('Faker Results:', fakerResults);
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
