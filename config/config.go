package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// 配置结构体
type Config struct {
	API struct {
		BaseURL string `yaml:"base_url"`
		Key     string `yaml:"key"`
	} `yaml:"api"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	OA struct {
		Url string `yaml:"url"`
	} `yaml:"oa"`
}

var GlobalConfig *Config

func MustInit() {
	GlobalConfig = loadConfig()
}

func loadConfig() *Config {
	filepath := composeConfigFileName("conf/config.yml", os.Getenv("SpecifiedConfig"))
	log.Printf("config filepath:%s", filepath)

	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var config Config
	if err = yaml.Unmarshal(f, &config); err != nil {
		panic(err)
	}
	return &config
}

func composeConfigFileName(basePath string, suffix string) string {
	var filepath = basePath

	if suffix != "" {
		filepath = strings.Join([]string{filepath, suffix}, ".")
	}

	return filepath
}
