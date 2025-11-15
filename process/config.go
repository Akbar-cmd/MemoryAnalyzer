package process

import "time"

type DisplayConfig struct {
	UpdateInterval time.Duration
	TopProcesses   int
}
