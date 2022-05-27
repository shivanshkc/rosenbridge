package configs

import (
	"sync"
)

var (
	// modelOnce ensures the singleton is instantiated only once.
	modelOnce = &sync.Once{}
	// modelSingleton points to the singleton value.
	modelSingleton *Model
)

// Get provides the config singleton.
func Get() *Model {
	modelOnce.Do(func() {
		modelSingleton = withViper()
	})
	return modelSingleton
}
