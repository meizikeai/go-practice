// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	App          App                        `mapstructure:"app"`
	MySQL        map[string][]MySQLInstance `mapstructure:"mysql"`
	Redis        map[string][]RedisInstance `mapstructure:"redis"`
	Kafka        map[string]KafkaInstance   `mapstructure:"kafka"`
	CryptoKey    CryptoKeyInstance          `mapstructure:"crypto"`
	JwtKey       JwtKeyInstance             `mapstructure:"jwt"`
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
	Addrs    []string `mapstructure:"addrs"`
	Password string   `mapstructure:"password"`
	DB       int      `mapstructure:"db"`
	PoolSize int      `mapstructure:"pool_size"`
}

type KafkaInstance struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
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

	path := "."
	if os.Getenv("GO_ENV") == "debug" {
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

	return result
}
