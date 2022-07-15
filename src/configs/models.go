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

	// Auth is the model of authentication configs.
	Auth struct {
		// InternalUsername is the username for internal basic auth.
		InternalUsername string `mapstructure:"internal_username"`
		// InternalPassword is the password for internal basic auth.
		InternalPassword string `mapstructure:"internal_password"`
	} `mapstructure:"auth"`

	// Bridges is the model of bridge-related configs.
	Bridges struct {
		// MaxBridgeLimit is the max number of bridges this node can host.
		MaxBridgeLimit int `mapstructure:"max_bridge_limit"`
		// MaxBridgeLimitPerClient is the max number of bridges this node can host per client.
		MaxBridgeLimitPerClient int `mapstructure:"max_bridge_limit_per_client"`
	} `mapstructure:"bridges"`

	// Discovery is the model of the discovery address related configs.
	Discovery struct {
		// DiscoveryAddr can be populated if it is known beforehand.
		DiscoveryAddr string `mapstructure:"discovery_addr"`
		// MaxAddrResolutionAttempts is max number of times we will attempt to resolve the discovery address.
		MaxAddrResolutionAttempts int `mapstructure:"max_addr_resolution_attempts"`
		// AddrResolutionPeriodSec is the number of seconds between two consecutive address resolution attempts.
		AddrResolutionPeriodSec int `mapstructure:"addr_resolution_period_sec"`
	} `mapstructure:"discovery"`

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

	// Mongo is the model of the MongoDB configs.
	Mongo struct {
		// Addr of the MongoDB deployment.
		Addr string `mapstructure:"addr"`
		// OperationTimeoutSec is the timeout in seconds for any MongoDB operation.
		OperationTimeoutSec int `mapstructure:"operation_timeout_sec"`
		// DatabaseName is the name of the logical database in MongoDB.
		DatabaseName string `mapstructure:"database_name"`
	} `mapstructure:"mongo"`
}
