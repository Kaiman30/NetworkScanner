package windows

import (
	"golang.org/x/sys/windows/registry"
)

// GetRegistryUint64 получает uint64 значение из реестра
func GetRegistryUint64(path, name string) (uint64, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.READ)
	if err != nil {
		return 0, err
	}
	defer key.Close()

	val, _, err := key.GetIntegerValue(name)
	return val, err
}

// GetRegistryString получает строковое значение из реестра
func GetRegistryString(path, name string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.READ)
	if err != nil {
		return "", err
	}
	defer key.Close()

	val, _, err := key.GetStringValue(name)
	return val, err
}

// GetRegistrySubKeys получает список подразделов реестра
func GetRegistrySubKeys(path string) ([]string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.READ)
	if err != nil {
		return nil, err
	}
	defer key.Close()

	return key.ReadSubKeyNames(-1)
}

// RegistryKeyExists проверяет существование ключа в реестре
func RegistryKeyExists(path string) bool {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.READ)
	if err != nil {
		return false
	}
	key.Close()
	return true
}
