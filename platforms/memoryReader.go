package platforms

import "MemoryAnalyzer/process"

type MemoryReader interface {
	// ReadSystemMemory Метод для чтения системной информации о памяти
	ReadSystemMemory() (process.SystemMemoryInfo, error)
	// GetProcessList Метод для получения списка процессов
	GetProcessList() ([]int, error)
	// ReadProcessMemory Метод чтения памяти конкретного процесса
	ReadProcessMemory(pid int) (uint64, error)
}
