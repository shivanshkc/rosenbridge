package logger

import (
	"sync"
)

var (
	// loggerOnce ensures the singleton is instantiated only once.
	loggerOnce = &sync.Once{}
	// loggerSingleton points to the singleton value.
	loggerSingleton Logger
)

// Get returns the Logger singleton.
func Get() Logger {
	// This statement only runs once.
	loggerOnce.Do(func() {
		loggerSingleton = newZapLogger()
	})
	// Returning the loaded configs.
	return loggerSingleton
}
