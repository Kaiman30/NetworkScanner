package checker

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/Kaiman30/NetworkChecker/pkg/windows"
)

type FakerCheckResult struct {
	CurrentConnection    *ConnectionInfo
	HostedNetworkActive  bool
	HostedNetworkSSID    string
	HostedNetworkClients int
	MobileHotspotActive  bool
	HotspotProfiles      []string
	VirtualAdapters      []AdapterInfo
	ConnectedDevices     []DeviceInfo
	SuspiciousActivities []string
	FakerDetected        bool
	FakerIndicators      []string
}

type ConnectionInfo struct {
	SSID              string
	State             string
	BSSID             string
	NetworkType       string
	RadioType         string
	Channel           string
	Signal            string
	IsHotspot         bool
	HotspotIndicators []string
}

type AdapterInfo struct {
	Name        string
	Description string
	MAC         string
}

type DeviceInfo struct {
	IP   string
	MAC  string
	Type string
}

// CheckFaker проверяет наличие признаков "фейкера"
func CheckFaker(ctx *CheckContext) *FakerCheckResult {
	result := &FakerCheckResult{
		SuspiciousActivities: []string{},
		FakerIndicators:      []string{},
		HotspotProfiles:      []string{},
	}

	// Проверяем текущее соединение
	result.CurrentConnection = getCurrentConnection()
	if result.CurrentConnection != nil && result.CurrentConnection.IsHotspot {
		result.SuspiciousActivities = append(result.SuspiciousActivities,
			fmt.Sprintf("Currently connected to hotspot: %s", result.CurrentConnection.SSID))
	}

	// Проверяем Windows Mobile Hotspot сервис
	result.MobileHotspotActive = isMobileHotspotRunning()
	if result.MobileHotspotActive {
		result.FakerDetected = true
		result.FakerIndicators = append(result.FakerIndicators,
			"Windows Mobile Hotspot service (icssvc) is running - FAKER INDICATOR")
		result.SuspiciousActivities = append(result.SuspiciousActivities,
			"Windows Mobile Hotspot service (icssvc) is running")
	}

	// Проверяем hosted network
	result.HostedNetworkActive, result.HostedNetworkSSID, result.HostedNetworkClients = getHostedNetworkStatus()
	if result.HostedNetworkActive {
		result.SuspiciousActivities = append(result.SuspiciousActivities,
			fmt.Sprintf("Active hosted network '%s' with %d client(s)", result.HostedNetworkSSID, result.HostedNetworkClients))
	}

	// Проверяем профили WiFi сетей
	result.HotspotProfiles = getHotspotProfiles()
	if len(result.HotspotProfiles) > 0 {
		result.SuspiciousActivities = append(result.SuspiciousActivities,
			fmt.Sprintf("Detected %d hotspot profile(s)", len(result.HotspotProfiles)))
	}

	// Проверяем виртуальные адаптеры
	result.VirtualAdapters = getVirtualAdapters()
	if len(result.VirtualAdapters) > 0 {
		result.SuspiciousActivities = append(result.SuspiciousActivities,
			fmt.Sprintf("%d virtual network adapter(s) detected", len(result.VirtualAdapters)))
	}

	// Проверяем подключенные устройства
	result.ConnectedDevices = getConnectedDevices()

	return result
}

func getCurrentConnection() *ConnectionInfo {
	output, err := windows.RunCommand("netsh", "wlan", "show", "interfaces")
	if err != nil {
		return nil
	}

	conn := &ConnectionInfo{
		HotspotIndicators: []string{},
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "SSID") && !strings.Contains(line, ":") {
			continue
		}
		if idx := strings.Index(line, "SSID"); idx != -1 {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				conn.SSID = strings.TrimSpace(parts[1])
			}
		}
		if idx := strings.Index(line, "State"); idx != -1 {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				conn.State = strings.TrimSpace(parts[1])
			}
		}
		if idx := strings.Index(line, "BSSID"); idx != -1 {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				conn.BSSID = strings.TrimSpace(parts[1])
			}
		}
		if idx := strings.Index(line, "Network type"); idx != -1 {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				conn.NetworkType = strings.TrimSpace(parts[1])
			}
		}
		if idx := strings.Index(line, "Radio type"); idx != -1 {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				conn.RadioType = strings.TrimSpace(parts[1])
			}
		}
		if idx := strings.Index(line, "Channel"); idx != -1 {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				conn.Channel = strings.TrimSpace(parts[1])
			}
		}
		if idx := strings.Index(line, "Signal"); idx != -1 {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				conn.Signal = strings.TrimSpace(parts[1])
			}
		}
	}

	if conn.State != "connected" || conn.SSID == "" {
		return nil
	}

	// Проверяем на признаки хотспота
	conn.IsHotspot, conn.HotspotIndicators = detectHotspotIndicators(conn)

	return conn
}

func detectHotspotIndicators(conn *ConnectionInfo) (bool, []string) {
	indicators := []string{}
	isHotspot := false

	// Проверка SSID на паттерны мобильных устройств
	hotspotPatterns := []string{
		"Android", "iPhone", "iPad", "Galaxy", "Pixel", "OnePlus",
		"Xiaomi", "Huawei", "Oppo", "Vivo", "Realme", "Nokia",
		"DIRECT-", "Redmi", "Mi ",
	}

	for _, pattern := range hotspotPatterns {
		if strings.Contains(conn.SSID, pattern) {
			isHotspot = true
			indicators = append(indicators, fmt.Sprintf("SSID matches mobile device pattern: %s", pattern))
			break
		}
	}

	// Проверка BSSID на MAC префиксы мобильных устройств
	if conn.BSSID != "" && conn.BSSID != "N/A" {
		oui := strings.ReplaceAll(strings.ToUpper(conn.BSSID[:8]), ":", "")

		mobileOUIs := map[string]string{
			"00505":  "Samsung",
			"0025BC": "Apple",
			"0026B":  "Apple",
			"A8667":  "Google Pixel",
			"F0D1A":  "Google",
			"5C8D4":  "Xiaomi",
			"F8A45":  "OnePlus",
			"DC44B":  "Huawei",
			"B0B98":  "Samsung Galaxy",
		}

		for prefix, device := range mobileOUIs {
			if strings.HasPrefix(oui, prefix) {
				isHotspot = true
				indicators = append(indicators, fmt.Sprintf("BSSID indicates %s device", device))
				break
			}
		}

		// Проверка на локально администрируемый адрес (второй символ 2, 6, A, E)
		if len(conn.BSSID) > 1 {
			secondChar := strings.ToUpper(string(conn.BSSID[1]))
			if strings.ContainsAny(secondChar, "26AE") {
				indicators = append(indicators, "BSSID uses locally administered address (common in hotspots)")
				isHotspot = true
			}
		}
	}

	// Проверка шлюза
	gateway := getDefaultGateway()
	if gateway != "" {
		if strings.HasPrefix(gateway, "192.168.137.") {
			isHotspot = true
			indicators = append(indicators, "Gateway indicates Windows PC Mobile Hotspot (192.168.137.x range) - FAKER INDICATOR")
		} else if strings.HasPrefix(gateway, "192.168.43.") {
			isHotspot = true
			indicators = append(indicators, "Gateway indicates Android hotspot (192.168.43.x range)")
		} else if gateway == "192.168.42.1" || gateway == "192.168.49.1" {
			isHotspot = true
			indicators = append(indicators, fmt.Sprintf("Gateway IP (%s) is typical for mobile hotspots", gateway))
		}

		// Проверка DNS
		dnsServers := getDNSServers()
		if len(dnsServers) > 0 && dnsServers[0] == gateway {
			indicators = append(indicators, "DNS server is same as gateway (typical hotspot configuration)")
			isHotspot = true
		}
	}

	// Проверка канала
	if conn.Channel != "" && conn.Channel != "N/A" {
		channelRegex := regexp.MustCompile(`^\d+`)
		if match := channelRegex.FindString(conn.Channel); match != "" {
			if match == "1" || match == "6" || match == "11" {
				indicators = append(indicators, fmt.Sprintf("Using common mobile hotspot channel: %s", match))
			}
		}
	}

	return isHotspot, indicators
}

func isMobileHotspotRunning() bool {
	// Проверяем статус сервиса icssvc (Mobile Hotspot)
	output, err := windows.RunCommand("sc", "query", "icssvc")
	if err != nil {
		return false
	}

	return strings.Contains(output, "RUNNING") || strings.Contains(output, "Running")
}

func getHostedNetworkStatus() (bool, string, int) {
	output, err := windows.RunCommand("netsh", "wlan", "show", "hostednetwork")
	if err != nil {
		return false, "N/A", 0
	}

	var status bool
	var ssid string
	var clients int

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "Status") {
			if strings.Contains(line, "Started") {
				status = true
			}
		}
		if strings.Contains(line, "SSID name") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				ssid = strings.Trim(strings.TrimSpace(parts[1]), `"`)
			}
		}
		if strings.Contains(line, "Number of clients") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &clients)
			}
		}
	}

	return status, ssid, clients
}

func getHotspotProfiles() []string {
	var profiles []string

	output, err := windows.RunCommand("netsh", "wlan", "show", "profiles")
	if err != nil {
		return profiles
	}

	hotspotPatterns := []string{
		"Android", "iPhone", "iPad", "Galaxy", "Pixel", "OnePlus",
		"Xiaomi", "DIRECT-", "SM-", "GT-",
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "All User Profile") {
			for _, pattern := range hotspotPatterns {
				if strings.Contains(line, pattern) {
					parts := strings.Split(line, ":")
					if len(parts) >= 2 {
						ssid := strings.TrimSpace(parts[1])
						if ssid != "" {
							profiles = append(profiles, ssid)
						}
					}
					break
				}
			}
		}
	}

	return profiles
}

func getVirtualAdapters() []AdapterInfo {
	var adapters []AdapterInfo

	// Используем WMI для получения информации о сетевых адаптерах
	// Для простоты, ищем через netsh
	output, err := windows.RunCommand("netsh", "interface", "show", "interface")
	if err != nil {
		return adapters
	}

	virtualPatterns := []string{"Virtual", "Hosted", "Wi-Fi Direct", "TAP"}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		for _, pattern := range virtualPatterns {
			if strings.Contains(line, pattern) {
				adapters = append(adapters, AdapterInfo{
					Description: line,
				})
				break
			}
		}
	}

	return adapters
}

func getConnectedDevices() []DeviceInfo {
	var devices []DeviceInfo

	output, err := windows.RunCommand("arp", "-a")
	if err != nil {
		return devices
	}

	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)\s+([0-9a-f-]+)\s+(dynamic|static)`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(strings.ToLower(line))
		if len(matches) >= 4 {
			ip := matches[1]
			mac := matches[2]
			devType := matches[3]

			// Фильтруем локальные адреса
			if strings.HasPrefix(ip, "192.168.") && devType == "dynamic" {
				devices = append(devices, DeviceInfo{
					IP:   ip,
					MAC:  mac,
					Type: devType,
				})
			}
		}
	}

	return devices
}

func getDefaultGateway() string {
	output, err := windows.RunCommand("ipconfig")
	if err != nil {
		return ""
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Default Gateway") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				gateway := strings.TrimSpace(parts[1])
				// Иногда может быть несколько шлюзов, берем первый
				if idx := strings.Index(gateway, ","); idx != -1 {
					gateway = gateway[:idx]
				}
				return strings.TrimSpace(gateway)
			}
		}
	}

	return ""
}

func getDNSServers() []string {
	var dnsServers []string

	output, err := windows.RunCommand("ipconfig", "/all")
	if err != nil {
		return dnsServers
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "DNS Servers") || strings.Contains(line, "DNS Server") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				dns := strings.TrimSpace(parts[1])
				if net.ParseIP(dns) != nil {
					dnsServers = append(dnsServers, dns)
				}
			}
		}
	}

	return dnsServers
}
