package memory

type SystemMemoryInfo struct {
	TotalMemory     uint64
	FreeMemory      uint64
	AvailableMemory uint64
	SwapTotal       uint64
	SwapFree        uint64
}

type ProcessInfo struct {
	// Идентификатор процесса
	PID int
	// Имя процесса
	Name string
	// Использование памяти в байтах
	MemoryUsage uint64
}
