package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type ProcessInfo struct {
	// Идентификатор процесса
	PID int
	// Имя процесса
	Name string
	// Использование памяти в байтах
	MemoryUsage uint64
}

func getShortProcessName(fullName string) string {

	// Получаем последний элемент пути
	file := filepath.Base(fullName)

	// Инициализируем суффиксы, которые нужно убрать, потом в цикле проверяем их наличие и удаляем
	suffixes := []string{" Helper (Renderer)", " Helper", "-helper (Renderer)"}
	for _, suf := range suffixes {
		if strings.HasSuffix(file, suf) {
			file = strings.TrimSuffix(file, suf)
		}
	}

	// Получаем расширение
	ext := filepath.Ext(file)
	// Удаляем расширение
	file = strings.TrimSuffix(file, ext)

	// Разбиваем строку по пробелам на части
	parts := strings.Fields(file)
	if len(parts) > 0 {
		// Записываем в file последний элемент
		file = parts[len(parts)-1]
	}

	const maxLen = 15
	// Считает кол-во символов
	if utf8.RuneCountInString(file) > maxLen {
		// Преобразуем строку в руны
		runes := []rune(file)
		// обрезаем до первых 12 символов + 3 точки
		file = string(runes[:maxLen-3]) + "..."
	}

	return file

}

func FormatTable(processes []ProcessInfo) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%-8s %-15s %10s\n", "PID", "NAME", "MEMORY")
	builder.WriteString("-----------------------------------\n")

	for _, proc := range processes {
		pid := fmt.Sprintf("%d", proc.PID)
		name := getShortProcessName(proc.Name)
		mem := FormatMemorySize(proc.MemoryUsage)

		fmt.Fprintf(&builder, "%-8s %-15s %10s\n", pid, name, mem)
	}

	return builder.String()
}

func GetProcessName(pid int) string {
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
