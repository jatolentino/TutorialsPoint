package constant

// constants for redis config
const (
	ClusterEnabled        = "redis.cluster.enabled"
	ClusterAddress        = "redis.cluster.node-addresses"
	SingleHost            = "redis.single.host"
	SinglePort            = "redis.single.port"
	IdleConnectionTimeout = "redis.idle-connection-timeout"
	MaxRedirects          = "redis.max-redirects"
	MinIdleConns          = "redis.min-idle-connections"
	ReadTimeout           = "redis.read-timeout"
	WriteTimeout          = "redis.write-timeout"
	MaxRetries            = "redis.max-retries"
	Timeout               = "redis.connect-timeout"
	PoolTimeout           = "redis.pool-timeout"
	PoolSize              = "redis.pool-size"
)