// internal/config/common.go
package config

import (
	"log"
	"os"
	"slices"

	"github.com/spf13/viper"
)

type Config struct {
	App          App                        `mapstructure:"app"`
	CryptoKey    CryptoKeyInstance          `mapstructure:"crypto"`
	LB           map[string]string          `mapstructure:"lb"`
	JwtKey       JwtKeyInstance             `mapstructure:"jwt"`
	Kafka        map[string]KafkaInstance   `mapstructure:"kafka"`
	MySQL        map[string][]MySQLInstance `mapstructure:"mysql"`
	Redis        map[string][]RedisInstance `mapstructure:"redis"`
	TencentCloud TencentCloudInstance       `mapstructure:"tencent_cloud"`
}

type App struct {
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
	Port string `mapstructure:"port"`
}

type MySQLInstance struct {
	Master          []string `mapstructure:"master"`
	Slave           []string `mapstructure:"slave"`
	MaxIdleConns    int      `mapstructure:"max_idle_conns"`
	MaxOpenConns    int      `mapstructure:"max_open_conns"`
	ConnMaxLifetime int      `mapstructure:"conn_max_lifetime"`
}

type RedisInstance struct {
	Addrs        []string `mapstructure:"addrs"`
	Password     string   `mapstructure:"password"`
	DB           int      `mapstructure:"db"`
	PoolSize     int      `mapstructure:"pool_size"`
	MinIdleConns int      `mapstructure:"min_idle_conns"`
}

type KafkaInstance struct {
	Brokers []string `mapstructure:"brokers"`
}

type CryptoKeyInstance struct {
	CurrentKey string `mapstructure:"current"` // Base64
	OldKey     string `mapstructure:"old"`     // Optional
}

type JwtKeyInstance struct {
	CurrentKeyID string `mapstructure:"current_key_id"` // "2025-01"
	CurrentKey   string `mapstructure:"current_key"`    // PEM
	OldKeyID     string `mapstructure:"old_key_id"`     // Optional
	OldKey       string `mapstructure:"old_key"`        // Optional
}

type TencentCloudInstance struct {
	SecretID  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
	Region    string `mapstructure:"region"`

	SMS struct {
		SdkAppID string `mapstructure:"sdk_app_id"`
		SignName string `mapstructure:"sign_name"`
	} `mapstructure:"sms"`

	SES struct {
		FromEmail string `mapstructure:"from_email"`
	} `mapstructure:"ses"`
}

func Load() *Config {
	var result *Config
	var env = []string{"release", "test"}
	var mode = os.Getenv("GO_ENV")
	var path = "."

	if !slices.Contains(env, mode) {
		mode = "test"
	}

	if mode == "release" {
		path = path + "/release"
	} else {
		path = path + "/test"
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("GO")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Read config failed: %v", err)
	}

	if err := viper.Unmarshal(&result); err != nil {
		log.Fatalf("Unmarshal config failed: %v", err)
	}

	result.App.Mode = mode

	return result
}
