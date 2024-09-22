package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	HTTPServer struct {
		IP   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"httpserver"`
}

func GetConfig(path string) *Config {
	config := new(Config)

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Config file error: " + err.Error())
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(config); err != nil {
		log.Fatal("Decoding config file error: " + err.Error())
	}
	log.Println("loaded config")
	return config

}
