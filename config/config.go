package config

import (
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    InfluxDBURL    string
    InfluxDBToken  string
    InfluxDBOrg    string
    InfluxDBBucket string
}

func LoadConfig() (*Config, error) {
    godotenv.Load() // Load .env file if it exists

    return &Config{
        InfluxDBURL:    getEnv("INFLUXDB_URL", "http://localhost:8086"),
        InfluxDBToken:  os.Getenv("INFLUXDB_TOKEN"),
        InfluxDBOrg:    getEnv("INFLUXDB_ORG", "myorg"),
        InfluxDBBucket: getEnv("INFLUXDB_BUCKET", "rules"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}