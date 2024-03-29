package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	configName = "configs"
	configType = "yaml"
)

// configsPaths is the list of locations that will be searched for the configs file.
var configPaths = []string{"/etc/", "/etc/rosenbridge/", "."}

// withViper loads the configs using spf13/viper.
// Panic is allowed here because configs are crucial to the application.
func withViper() *Model {
	// Specifying the configs file name and type to viper.
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	// Adding configs paths to viper.
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	// Reading the configs file.
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error in ReadInConfig: %w", err))
	}

	model := &Model{}
	// Unmarshalling it to the model instance.
	if err := viper.Unmarshal(model); err != nil {
		panic(fmt.Errorf("error in Unmarshal: %w", err))
	}

	return model
}
