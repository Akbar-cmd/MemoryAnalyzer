package process

import (
	"fmt"
	"sort"
	"time"
)

func DisplayDashboard(stats SystemMemoryInfo, processes []ProcessInfo, config DisplayConfig) {
	// Очистка экрана и перемещение курсора
	fmt.Print("\033[H\033[2J")
	fmt.Print("=== Memory Analyzer ===\n\n")

	fmt.Println(FormatSystemStats(stats))

	fmt.Print("Top Memory Processes:\n")
	//Сортировка процессов по памяти по убыванию
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].MemoryUsage > processes[j].MemoryUsage
	})
	// Ограничение кол-ва процессов
	if len(processes) > config.TopProcesses {
		processes = processes[:config.TopProcesses]
	}
	fmt.Println(FormatTable(processes))

	currTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("Updated: %s\n", currTime)
	fmt.Println("Press Ctrl+C to exit")
}
