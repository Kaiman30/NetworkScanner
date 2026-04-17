package web

// GetHTMLTemplate возвращает HTML шаблон для отчета
func GetHTMLTemplate() string {
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
            • "Настроек не выявлено" - параметры в оптимальном состоянии<br>
            • "Обнаружены настройки" - параметры изменены от стандартных
        </div>
        <footer>
            Network Settings Checker by Kaiman4ik :3
        </footer>
    </div>
    <script>
        // Данные передаются из Go
        const rawData = document.currentScript.getAttribute('data-results');
        let results;
        
        try {
            results = JSON.parse(rawData || '{}');
        } catch(e) {
            console.error('Parse error:', e);
            results = { passed: [], failed: [] };
        }
        
        // Убеждаемся, что массивы существуют
        results.passed = results.passed || [];
        results.failed = results.failed || [];
        
        // Обновляем счетчики
        document.getElementById('passedCount').textContent = results.passed.length;
        document.getElementById('failedCount').textContent = results.failed.length;
        
        // Функция отображения списка
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
        
        // Отображаем списки
        renderList('passedList', results.passed, '✨ Все параметры в оптимальном состоянии');
        renderList('failedList', results.failed, '🔧 Нет измененных параметров');
        
        // Время отчета
        document.getElementById('timestamp').textContent = new Date().toLocaleString();
        
        console.log('Report loaded:', results);
    </script>
</body>
</html>`
}
