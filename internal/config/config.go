package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// структура под будущую бд
type PgConf struct {
}

type AppConf struct {
	Host     string
	Port     string
	Mode     string
	LogLevel string
}
type JWTConf struct {
	Secret string
}
type Log struct {
	Level         string
	Encoding      string
	FilePath      string
	EnableConsole bool
	EnableFile    bool
}
type SenderConf struct {
	Endpoint string
}
type Config struct {
	Pg     PgConf
	JWT    JWTConf
	App    AppConf
	Log    Log
	Sender SenderConf
}

func init() {
	println("Loading .env file...")
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: Could not load .env file")
		return
	}
	println("Successfully Loaded .env file")
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func New() *Config {
	return &Config{
		Pg: PgConf{},
		App: AppConf{
			Host:     getEnv("HTTP_HOST", "0.0.0.0"),
			Port:     getEnv("HTTP_PORT", "8080"),
			Mode:     getEnv("APP_MODE", "debug"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Log: Log{
			Level:         getEnv("LOG_LEVEL", "info"),
			Encoding:      getEnv("LOG_ENCODING", "console"),
			FilePath:      getEnv("LOG_FILE_PATH", "provider.log"),
			EnableConsole: getEnv("LOG_ENABLE_CONSOLE", "true") == "true",
			EnableFile:    getEnv("LOG_ENABLE_FILE", "true") == "true",
		},
		Sender: SenderConf{
			Endpoint: getEnv("SENDER_ENDPOINT", "http://localhost:8081/users"),
		},
		JWT: JWTConf{
			Secret: getEnv("JWT_SECRET", "secret"),
		},
	}
}
