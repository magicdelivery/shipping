package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type App struct {
	ListenAddr string	`yaml:"listen_addr"`
}

type RedisStorage struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ParcelLockerService struct {
	EndpointUrl string `yaml:"endpoint_url"`
	CacheTtl    int    `yaml:"cache_ttl"`
}

type Config struct {
	App					App					`yaml:"app"`
	RedisStorage        RedisStorage		`yaml:"redis_storage"`
	ParcelLockerService ParcelLockerService	`yaml:"parcel_locker_service"`
}

func NewConfig(configPath string) (*Config, error) {
	config, err := loadFromYaml(&Config{}, configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load config file \"%s\", %s", configPath, err)
	}
	config = loadFromEnv(config)
	log.Printf("%v", config)
	return config, nil
}

func loadFromYaml(config *Config, configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func loadFromEnv(config *Config) *Config {
	config.App.ListenAddr = GetStrEnv("APP_LISTEN_ADDR", config.App.ListenAddr)
	config.RedisStorage.Host = GetStrEnv("REDIS_HOST", config.RedisStorage.Host)
	config.RedisStorage.Port = GetIntEnv("REDIS_PORT", config.RedisStorage.Port)
	config.ParcelLockerService.EndpointUrl = GetStrEnv("PARCEL_LOCKER_SERVICE_ENDPOINT_URL", config.ParcelLockerService.EndpointUrl)
	config.ParcelLockerService.CacheTtl = GetIntEnv("PARCEL_LOCKER_CACHE_TTL", config.ParcelLockerService.CacheTtl)
	return config
}