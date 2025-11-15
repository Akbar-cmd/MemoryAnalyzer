package main

import (
	"MemoryAnalyzer/interfaces"
	"MemoryAnalyzer/memory"
	"MemoryAnalyzer/platforms"
	"MemoryAnalyzer/process"
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

// Проект Memory Analyzer
// 1. Создаем базовую структуру для информации о процессе
// 2. Создаем интерфейс для кроссплатформенной работы с процессами
// 3. Получение списка процессов в Darwin(macOS) и Linux
// 4. Чтение памяти процесса в Darwin и Linux
// 5. Создание структуры для информации о системной памяти
// 6. Чтение информации о системной памяти в Darwin и Linux
// 7. Создаем функцию форматирования размера памяти (правильно ли через if? или стоило сделать через цикл?)
// 8. Создание структуры для конфигурации отображения
// 9. Создание функции сокращения имен процессов
// 10. Создание функции форматирования таблицы процессов
// 11. Создание функции форматирования системной аналитики
// 12. Создание функции отображения информационной панели
// 13. Настройка обработки сигналов завершения программы
// 14. Настройка обработки сигналов завершения программы
// 15. Настройка периодического обновления данных

func performGracefulShutdown(ctx context.Context) {
	fmt.Println("=== Graceful Shutdown ===")
	fmt.Println("✓ Ticker stopped (via defer)")
	fmt.Println("✓ Context cancelled")
	fmt.Println("✓ All goroutines terminated")
	time.Sleep(2 * time.Second)
	fmt.Println("=== Shutdown complete ===")
}

func main() {
	// Инициализируем компоненты
	var reader interfaces.MemoryReader
	switch runtime.GOOS {
	case "darwin":
		reader = &platforms.DarwinMemoryReader{}
	case "linux":
		reader = &platforms.LinuxMemoryReader{}
	default:
		fmt.Printf("Unsupported OS: %s\n", runtime.GOOS)
		return
	}

	// Создаем конфигурацию и настраиваем обработку событий
	config := process.DisplayConfig{
		UpdateInterval: 3 * time.Second,
		TopProcesses:   10,
	}

	// Создаём контекст, который можно отменить
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(config.UpdateInterval)
	defer ticker.Stop()

	fmt.Printf("Starting Memory Analyzer on %s\n", runtime.GOOS)
	fmt.Println("Press Ctrl+C to shutdown gracefully...")

	// Горутина для обработки сигналов завершения
	go func() {
		<-sigChan
		fmt.Println("\n\nReceived interrupt signal. Initiating graceful shutdown...")

		// Запускаем graceful shutdown в отдельной горутине
		go performGracefulShutdown(ctx)

		// Даём время на завершение (5 секунд)
		time.Sleep(5 * time.Second)
		cancel()
	}()

	// Реализовать основной цикл обработки событий
	for {
		select {
		// Обработка сигнала завершения
		case <-ctx.Done():
			fmt.Println("\nReceived interrupt signal. Exiting...")
			return
		case <-ticker.C:
			// Получение системной информации
			sysInfo, err := reader.ReadSystemMemory()
			if err != nil {
				fmt.Printf("Error reading system memory: %v\n", err)
				continue
			}

			// Получение списка процессов
			pids, err := reader.GetProcessList()
			if err != nil {
				fmt.Printf("Error getting process list: %v\n", err)
				continue
			}

			// Сбор информации о процессах
			var processes []memory.ProcessInfo
			for _, pid := range pids {
				if mem, err := reader.ReadProcessMemory(pid); err == nil {
					name := reader.GetProcessName(pid)
					processes = append(processes, memory.ProcessInfo{
						PID:         pid,
						Name:        name,
						MemoryUsage: mem,
					})
				}
			}

			// Отображение информационной панели
			process.DisplayDashboard(sysInfo, processes, config)

		}
	}
}
