package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config encapsulates all config required by the application.
type Config struct {
	HttpServer struct {
		Addr           string   `json:"addr"`
		AllowedOrigins []string `json:"allowedOrigins"`
		// Read here: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Max-Age
		CorsMaxAgeSec int `json:"corsMaxAgeSec"`
	} `json:"httpServer"`

	Logger struct {
		Level  string `json:"level"`
		Pretty bool   `json:"pretty"`
	} `json:"logger"`

	Database struct {
		UsersFilePath string `json:"usersFilePath"`
	} `json:"database"`

	Frontend struct {
		// The base URL of the backend that the frontend will use.
		BackendAddr string `json:"backendAddr"`
		// Path to the directory that contains the frontend SPA files.
		Path string `json:"path"`
	} `json:"frontend"`
}

// Load config from the given JSON file.
func Load(jsonPath string) (Config, error) {
	content, err := os.ReadFile(jsonPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file at %s because: %w", jsonPath, err)
	}

	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config file at %s because: %w", jsonPath, err)
	}

	if err := setFrontendConfig(config); err != nil {
		return Config{}, fmt.Errorf("failed to set frontend config because: %w", err)
	}

	return config, nil
}

// setFrontendConfig sets the config values for the frontend.
func setFrontendConfig(conf Config) error {
	if conf.Frontend.Path == "" {
		return nil
	}

	path := filepath.Join(conf.Frontend.Path, "config.json")
	data := fmt.Sprintf(`{ "apiBaseUrl": "%s" }`, conf.Frontend.BackendAddr)

	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		return fmt.Errorf("error writing frontend config file: %w", err)
	}

	return nil
}
