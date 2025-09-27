package debug

import (
	"time"
)

func Timer(name string) func() {
	//log.Debug("Starting timer", "name", name)
	start := time.Now()
	return func() {
		duration := time.Since(start)
		if duration > 0 {
			log.Debug("Ending timer", "name", name, "duration", duration)
		}
	}
}
