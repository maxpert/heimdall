package main

import (
	"github.com/spf13/viper"
	"gopkg.in/inconshreveable/log15.v2"
	"heimdall/heimdall"
	"os"
)

func initConfig() (*heimdall.Config, error) {
	viperCfg := viper.New()
	viperCfg.SetDefault("max_workers", 64)
	viperCfg.SetDefault("database.host", "127.0.0.1")
	viperCfg.SetDefault("database.port", 5432)
	viperCfg.SetDefault("database.user", "root")
	viperCfg.SetDefault("database.password", "")
	viperCfg.SetDefault("database.name", nil)
	viperCfg.SetDefault("database.max_connections", 16)
	viperCfg.SetDefault("database.procedure_name", "notify_heimdall")

	viperCfg.SetConfigName("config")
	viperCfg.SetConfigType("yaml")
	viperCfg.AddConfigPath("/etc/heimdall/")
	viperCfg.AddConfigPath("$HOME/.heimdall")
	viperCfg.AddConfigPath(".")

	if err := viperCfg.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &heimdall.Config{}
	if err := viperCfg.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func main() {
	logger := log15.New("MODULE", "main")
	cfg, err := initConfig()
	if err != nil {
		logger.Crit("unable to unmarshal configuration", "error", err)
		os.Exit(1)
	}

	if cfg.DB.Host == "" {
		logger.Crit("invalid host configuration", "config", cfg)
		os.Exit(1)
	}

	connectionPool, err := heimdall.Bootup(cfg)
	if err != nil {
		logger.Crit(err.Error())
		os.Exit(1)
	}

	err = heimdall.Listen(connectionPool, cfg)
	if err != nil {
		logger.Crit(err.Error())
		os.Exit(1)
	}
}
