package echolet

import (
	"encoding/json"
	"log"
	"os"
)

type ServerConfig struct {
	Name         string `json:"name"`
	Developer    string `json:"developer"`
	Version      string `json:"version"`
	Branch       string `json:"branch"`
	Port         int    `json:"port"`
	ContextPath  string `json:"context_path"`
	LogLevel     string `json:"log_level"`
	EnableLogger bool   `json:"enable_logger"`
}

type MySqlDSN struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	DbName   string `json:"dbname"`
	PoolSize int    `json:"pool_size"`
}

func LoadMySqlDSN() *MySqlDSN {
	return loadConfigFile("config/mysqlc.json", &MySqlDSN{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Pass:     "pass4pass",
		DbName:   "gorux",
		PoolSize: 10,
	}).(*MySqlDSN)
}

// LoadServerCofing...
func LoadServerConfig() *ServerConfig {
	return loadConfigFile("config/server.json", &ServerConfig{
		Name:         "Biz WebAPI Server",
		Developer:    "Steven Chen",
		Version:      "v1.0.0",
		Branch:       "Dev",
		LogLevel:     "off",
		EnableLogger: false,
	}).(*ServerConfig)
}

func loadConfigFile(file string, conf interface{}) interface{} {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Can't open config file: ", err)
	}
	defer f.Close()
	jsonParser := json.NewDecoder(f)
	err = jsonParser.Decode(conf)
	if err != nil {
		log.Fatal("Parse config failed: ", err)
	}
	return conf
}
