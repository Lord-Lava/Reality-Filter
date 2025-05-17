package ports

// ConfigProvider defines the interface for accessing application configuration
type ConfigProvider interface {
	GetMongoDBConfig() MongoDBConfig
	GetRedisConfig() RedisConfig
	GetPostgresConfig() PostgresConfig
	GetLogConfig() LogConfig
}

// MongoDBConfig represents MongoDB configuration requirements
type MongoDBConfig interface {
	GetURI() string
	GetDatabase() string
}

// RedisConfig represents Redis configuration requirements
type RedisConfig interface {
	GetAddr() string
	GetPassword() string
	GetDB() int
}

// PostgresConfig represents PostgreSQL configuration requirements
type PostgresConfig interface {
	GetDSN() string
}

// LogConfig represents logging configuration requirements
type LogConfig interface {
	GetLevel() string
	GetFormat() string
	GetOutputPath() string
}
