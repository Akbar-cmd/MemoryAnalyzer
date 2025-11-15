package process

import (
	"MemoryAnalyzer/memory"
	"fmt"
	"strings"
)

func FormatMemorySize(bytes uint64) string {

	// объявляем единицы измерения в виде const
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}

}

func FormatSystemStats(stats memory.SystemMemoryInfo) string {
	// создаем builder для строки
	var builder strings.Builder
	// для хранения %
	var usedPerc float64
	var swapPerc float64

	// Вызываем FormatMemorySize() для преобразования значений памяти
	total := FormatMemorySize(stats.TotalMemory)
	// Вычисляем использованную память в байтах
	usedBytes := stats.TotalMemory - stats.FreeMemory
	// Проверяем равна ли общая память 0
	if stats.TotalMemory > 0 {
		usedPerc = float64(usedBytes) / float64(stats.TotalMemory) * 100
	}
	// Форматируем процент с одним знаком после запятой
	usedForm := fmt.Sprintf("%.1f%%", usedPerc)
	// Переводим в строку
	usedMem := fmt.Sprintf("%s (%s)", FormatMemorySize(usedBytes), usedForm)

	available := FormatMemorySize(stats.AvailableMemory)

	swapBytes := stats.SwapTotal - stats.SwapFree
	if stats.SwapTotal > 0 {
		swapPerc = float64(swapBytes) / float64(stats.SwapTotal) * 100
	}
	swapForm := fmt.Sprintf("%.1f%%", swapPerc)
	swapUsedMem := fmt.Sprintf("%s (%s)", FormatMemorySize(swapBytes), swapForm)

	// Делаем итоговый вывод
	fmt.Fprintf(&builder,
		"System Memory:\n"+
			"Total:     %s\n"+
			"Used:      %s\n"+
			"Available:  %s\n"+
			"Swap Used:  %s\n",
		total,
		usedMem,
		available,
		swapUsedMem,
	)

	return builder.String()

}
