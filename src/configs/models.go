package configs

// Model represents the configs model.
type Model struct {
	// Application is the model of application configs.
	Application struct {
		// Name of the application.
		Name string `mapstructure:"name"`
		// Version of the application.
		Version string `mapstructure:"version"`
	} `mapstructure:"application"`

	// HTTPServer is the model of the HTTP Server configs.
	HTTPServer struct {
		// Addr is the address of the HTTP server.
		Addr string `mapstructure:"addr"`
	} `mapstructure:"http_server"`

	// Logger is the model of the logger configs.
	Logger struct {
		// Level for logging.
		Level string `mapstructure:"level"`
	} `mapstructure:"logger"`
}
