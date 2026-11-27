package config

type MergedConfig struct {
	Name         string
	ID           string
	Container    string
	Port         string
	JWT          JWTConfig
	Database     DatabaseConfig
	Redis        RedisConfig
	Dependencies map[string]DependencyConfig
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type DependencyConfig struct {
	Name string
	ID   string
}

type JWTConfig struct {
	ServiceSecret []byte
	GatewaySecret []byte
	LinkSecret    []byte
}
