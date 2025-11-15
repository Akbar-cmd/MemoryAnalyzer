package platforms

import (
	"MemoryAnalyzer/process"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// Тестовая реализация для проверки интерфейса
type DarwinMemoryReader struct{}

func (d *DarwinMemoryReader) GetProcessList() ([]int, error) {
	// Создаем команду для выполнения
	cmd := exec.Command("ps", "-e", "-o", "pid=")

	// Исполняем команду и получаем вывод
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error output: %s\n", string(output))
		return []int{}, err
	}

	// Разбиваем строки при переходе на следующую строку
	strOut := strings.Split(string(output), "\n")

	// Создаем слайс для преобразованных в int процессов
	pids := make([]int, 0)

	// Проходимся по процессам, преобразуем каждый в int и записываем в слайс
	for _, pid := range strOut {
		// удаляем пробелы, иначе не сможет конвертировать в int
		pid = strings.TrimSpace(pid)

		// пропускаем пустую строку, которая появляется после strings.Splint("\n")
		if pid == "" {
			continue
		}

		// конвертируем строку в int
		intOut, err := strconv.Atoi(pid)
		if err != nil {
			log.Println("Не удалось перевести string в int", pid)
			continue
		}
		// Добавляем в слайс
		pids = append(pids, intOut)
	}

	return pids, nil
}

func (d *DarwinMemoryReader) ReadProcessMemory(pid int) (uint64, error) {

	// Создаем команду для выполнения и тк rss требует указания конкретного pid через флаг, то надо посредством Sprintf вставить в параметры
	cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", pid), "-o", "rss=")

	// Получаем вывод
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error output: %s\n", string(output))
		return 0, err
	}

	// Тк здесь могут быть пробелы в начале числа, например  " 1920", поэтому удаляем пробелы, чтобы преобразовать в Uint
	str := strings.TrimSpace(string(output))

	// Конвертируем в Uint
	// rss возвращает в килобайтах
	rssKb, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		fmt.Println("Проблема конвертации string в uint", str)
		return 0, err
	}

	// Переводим в байты
	rssByte := rssKb * 1024

	return rssByte, nil
}

func (d *DarwinMemoryReader) ReadSystemMemory() (process.SystemMemoryInfo, error) {
	info := process.SystemMemoryInfo{}

	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Combined Output:\n%s\n", string(output))
		return info, err
	}

	totalMemStr := strings.TrimSpace(string(output))
	totalMem, err := strconv.ParseUint(totalMemStr, 10, 64)
	if err != nil {
		return info, err
	}

	info.TotalMemory = totalMem

	cmd = exec.Command("vm_stat")
	output, err = cmd.Output()
	if err != nil {
		return info, err
	}

	vmLines := strings.Split(string(output), "\n")
	freePages := uint64(0)
	inactivePages := uint64(0)

	for _, line := range vmLines {
		if strings.Contains(line, "Pages free:") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				lastPart := parts[len(parts)-1]
				lastPart = strings.TrimSuffix(lastPart, ".")
				freePages, err = strconv.ParseUint(lastPart, 10, 64)
				if err != nil {
					fmt.Printf("Error parsing free pages: %v\n", err)
					continue
				}
			}
		}
		if strings.Contains(line, "Pages inactive:") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				lastPart := parts[len(parts)-1]
				lastPart = strings.TrimSuffix(lastPart, ".")
				inactivePages, err = strconv.ParseUint(lastPart, 10, 64)
				if err != nil {
					fmt.Printf("Error parsing inactive pages: %v\n", err)
					continue
				}
			}
		}
	}

	pageSize := uint64(4096)
	info.FreeMemory = freePages * pageSize
	info.AvailableMemory = (freePages + inactivePages) * pageSize

	cmd = exec.Command("sysctl", "-n", "vm.swapusage")
	output, err = cmd.Output()
	if err != nil {
		return info, err
	}

	outputStr := strings.TrimSpace(string(output))

	if strings.Contains(outputStr, "total") {
		parts := strings.Split(outputStr, " ")
		for i, part := range parts {
			if part == "total" && i+2 < len(parts) {
				totalSwapStr := parts[i+2]
				info.SwapTotal = parseSwap(totalSwapStr)
			}
			if part == "free" && i+2 < len(parts) {
				freeSwapStr := parts[i+2]
				info.SwapFree = parseSwap(freeSwapStr)
			}
		}
	}

	return info, nil
}

func parseSwap(value string) uint64 {
	mltplr := uint64(1)

	if strings.HasSuffix(value, "G") {
		mltplr = 1024 * 1024 * 1024
		value = strings.TrimSuffix(value, "G")
	} else if strings.HasSuffix(value, "M") {
		mltplr = 1024 * 1024
		value = strings.TrimSuffix(value, "M")
	} else if strings.HasSuffix(value, "K") {
		mltplr = 1024
		value = strings.TrimSuffix(value, "K")
	}

	floatVal, _ := strconv.ParseFloat(value, 64)
	return uint64(floatVal * float64(mltplr))
}
