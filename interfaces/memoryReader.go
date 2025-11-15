package interfaces

import "MemoryAnalyzer/memory"

type MemoryReader interface {
	// ReadSystemMemory Метод для чтения системной информации о памяти
	ReadSystemMemory() (memory.SystemMemoryInfo, error)
	// GetProcessList Метод для получения списка процессов
	GetProcessList() ([]int, error)
	// ReadProcessMemory Метод чтения памяти конкретного процесса
	ReadProcessMemory(pid int) (uint64, error)
	GetProcessName(pid int) string
}
