package platforms

import (
	"MemoryAnalyzer/memory"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LinuxMemoryReader struct{}

func (l *LinuxMemoryReader) GetProcessList() ([]int, error) {

	// Читаем директорию
	dir, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Printf("Error output: %s\n", dir)
		return []int{}, err
	}

	// Создаем слайс интов
	pids := make([]int, 0)

	// Читаем каждый элемент директории
	for _, pid := range dir {
		// Проверяем, что это директория
		if !pid.IsDir() {
			continue
		}

		// тк pid выглядит подобно "d acpi/", то вызываем pid.Name(), чтоб получить лишь имя процессов
		name := pid.Name()
		// Конвертируем в число
		intPid, err := strconv.Atoi(name)
		// Тк не все pid равны цифрам, то если оно таковым не является - пропускаем
		if err != nil {
			continue
		}

		// добавляем в слайс
		pids = append(pids, intPid)
	}

	return pids, nil
}

func (l *LinuxMemoryReader) ReadProcessMemory(pid int) (uint64, error) {
	// Этап 1: Формируем путь к файлу
	filePath := fmt.Sprintf("/proc/%d/status", pid)

	// Этап 2: Читаем файл
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	// Этап 3: Разбиваем на строки
	lines := strings.Split(string(content), "\n")

	// Этап 4 & 5: Ищем VmRSS и извлекаем число
	for _, line := range lines {
		if strings.HasPrefix(line, "VmRSS:") {
			// Удаляем префикс "VmRSS:"
			trimmed := strings.TrimPrefix(line, "VmRSS:")
			// Удаляем пробелы
			trimmed = strings.TrimSpace(trimmed)
			// Удаляем суффикс " kB"
			trimmed = strings.TrimSuffix(trimmed, " kB")

			// Преобразуем в число
			vmRSSKB, err := strconv.ParseUint(trimmed, 10, 64)
			if err != nil {
				return 0, err
			}

			// Этап 6: Конвертируем килобайты в байты
			vmRSSBytes := vmRSSKB * 1024
			return vmRSSBytes, nil
		}
	}

	// Если VmRSS не найдена
	return 0, fmt.Errorf("VmRSS not found in /proc/%d/status", pid)
}

func (l *LinuxMemoryReader) ReadSystemMemory() (memory.SystemMemoryInfo, error) {
	info := memory.SystemMemoryInfo{}

	file, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return info, err
	}

	lines := strings.Split(string(file), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				memTotalKB, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return info, err
				}
				info.TotalMemory = memTotalKB * 1024
			}
		}

		if strings.HasPrefix(line, "MemFree") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				memFreeKB, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return info, err
				}
				info.FreeMemory = memFreeKB * 1024
			}
		}

		if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				memAvailKB, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return info, err
				}
				info.AvailableMemory = memAvailKB * 1024
			}
		}

		if strings.HasPrefix(line, "SwapTotal:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				swapTotalKB, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return info, err
				}
				info.SwapTotal = swapTotalKB * 1024
			}
		}

		if strings.HasPrefix(line, "SwapFree:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				swapFreeKB, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return info, err
				}
				info.SwapFree = swapFreeKB * 1024
			}
		}
	}

	return info, nil
}

func (l *LinuxMemoryReader) GetProcessName(pid int) string {
	filePath := fmt.Sprintf("/proc/%d/comm", pid)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("%d", pid) // Fallback на PID если ошибка
	}
	name := strings.TrimSpace(string(content))
	if name == "" {
		return fmt.Sprintf("%d", pid)
	}
	return name
}
