package config

import "github.com/spf13/viper"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	CoinDesk CoinDeskConfig `mapstructure:"coindesk"`
	Database DatabaseConfig `mapstructure:"database"`
	Cron     CronConfig     `mapstructure:"cron"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type CoinDeskConfig struct {
	URL        string `mapstructure:"url"`
	APIKey     string `mapstructure:"api_key"`
	TimeoutSec int    `mapstructure:"timeout_sec"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"`
}

type CronConfig struct {
	IntervalMin int `mapstructure:"interval_min"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
