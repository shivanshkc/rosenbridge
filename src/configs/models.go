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
		ClusterUsername string `mapstructure:"cluster_username"`
		ClusterPassword string `mapstructure:"cluster_password"`
	} `mapstructure:"auth"`

	// HTTPServer is the model of the HTTP Server configs.
	HTTPServer struct {
		// Addr is the address of the HTTP server.
		Addr string `mapstructure:"addr"`
		// DiscoveryAddr is the address of this node that other nodes can use to reach it.
		DiscoveryAddr string `mapstructure:"discovery_addr"`
		// DiscoveryProtocol is the protocol to be used while contacting a cluster node (http or https).
		DiscoveryProtocol string `mapstructure:"discovery_protocol"`
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
