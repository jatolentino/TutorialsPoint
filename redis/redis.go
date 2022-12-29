package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	goredis "github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"github.com/thejasn/go-redis-template/pkg/redis/config"
	"github.com/thejasn/go-redis-template/pkg/redis/constant"
)

type yamlConfig struct {
	NodeAddresses         []string
	Enabled               bool
	Host                  string
	Port                  string
	IdleConnectionTimeout time.Duration
	ConnectTimeout        time.Duration
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	MaxRetries            int
	MaxRedirects          int
	MinIdleConns          int
	PoolSize              int
	PoolTimeout           time.Duration
}

// Config for connecting to a redis instance
type Config struct {
	SingleConfig  *goredis.Options
	ClusterConfig *goredis.ClusterOptions
}

// Client represents either a redis cluster client or a standalone client
type Client struct {
	Cluster *goredis.ClusterClient
	Single  *goredis.Client
}

// BuildConfig builds the redis configuration parameters from the application.yaml
func buildConfig() (Config, bool) {

	redisConfig := viper.New()
	config.LoadConfig(redisConfig)
	redisConfig.WatchConfig()
	redisConfig.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		config.LoadConfig(redisConfig)
	})

	yConfig := yamlConfig{
		NodeAddresses:         strings.Split(redisConfig.GetString(constant.ClusterAddress), ","),
		Enabled:               redisConfig.GetBool(constant.ClusterEnabled),
		Host:                  redisConfig.GetString(constant.SingleHost),
		Port:                  redisConfig.GetString(constant.SinglePort),
		MaxRetries:            redisConfig.GetInt(constant.MaxRetries),
		MaxRedirects:          redisConfig.GetInt(constant.MaxRedirects),
		MinIdleConns:          redisConfig.GetInt(constant.MinIdleConns),
		ReadTimeout:           time.Duration(redisConfig.GetInt(constant.ReadTimeout)) * time.Second,
		WriteTimeout:          time.Duration(redisConfig.GetInt(constant.WriteTimeout)) * time.Second,
		IdleConnectionTimeout: time.Duration(redisConfig.GetInt(constant.IdleConnectionTimeout)) * time.Second,
		ConnectTimeout:        time.Duration(redisConfig.GetInt(constant.Timeout)) * time.Second,
		PoolTimeout:           time.Duration(redisConfig.GetInt(constant.PoolTimeout)) * time.Second,
		PoolSize:              redisConfig.GetInt(constant.PoolSize),
	}

	if yConfig.Enabled {
		return Config{
			ClusterConfig: &goredis.ClusterOptions{
				Addrs:        yConfig.NodeAddresses,
				MaxRedirects: yConfig.MaxRedirects,
				MaxRetries:   yConfig.MaxRetries,
				DialTimeout:  yConfig.ConnectTimeout,
				ReadTimeout:  yConfig.ReadTimeout,
				WriteTimeout: yConfig.WriteTimeout,
				PoolSize:     yConfig.PoolSize,
				MinIdleConns: yConfig.MinIdleConns,
				IdleTimeout:  yConfig.IdleConnectionTimeout,
				PoolTimeout:  yConfig.PoolTimeout,
			},
		}, true
	}
	return Config{
		SingleConfig: &goredis.Options{
			Addr:         strings.Join([]string{yConfig.Host, yConfig.Port}, ":"),
			MaxRetries:   yConfig.MaxRetries,
			DialTimeout:  yConfig.ConnectTimeout,
			ReadTimeout:  yConfig.ReadTimeout,
			WriteTimeout: yConfig.WriteTimeout,
			PoolSize:     yConfig.PoolSize,
			PoolTimeout:  yConfig.PoolTimeout,
			MinIdleConns: yConfig.MinIdleConns,
			IdleTimeout:  yConfig.IdleConnectionTimeout,
		},
	}, false

}

// NewClient generates a new redis client
func NewClient() Client {
	redisConfig, isCluster := buildConfig()
	if isCluster {
		return Client{
			Cluster: goredis.NewClusterClient(redisConfig.ClusterConfig),
		}
	}
	return Client{
		Single: goredis.NewClient(redisConfig.SingleConfig),
	}
}