package config

import (
	"gopkg.in/yaml.v2"
	"log/slog"
	"os"
)

type Config struct {
	HTTPServer struct {
		IP               string   `yaml:"ip"`
		Port             string   `yaml:"port"`
		AllowedIPsByCORS []string `yaml:"allowed_ips_by_cors"`
	} `yaml:"httpserver"`
	Postgres struct {
		IP       string `yaml:"ip"`
		Port     string `yaml:"port"`
		DBname   string `yaml:"dbname"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"postgres"`
	Redis struct {
		IP       string `yaml:"ip"`
		Port     string `yaml:"port"`
		DBnum    int    `yaml:"dbnum"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
}

func GetConfig(path string) (*Config, error) {
	config := new(Config)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err = d.Decode(config); err != nil {
		return nil, err
	}
	slog.Info("loaded config")
	return config, nil

}
