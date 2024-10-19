package main

import (
    "log"
    "os"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/influxdata/influxdb-client-go/v2"
    "rule-engine/api"
    "rule-engine/config"
    "rule-engine/engine"
    "rule-engine/storage"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize InfluxDB client
    client := influxdb2.NewClient(cfg.InfluxDBURL, cfg.InfluxDBToken)
    defer client.Close()

    // Initialize storage
    store := storage.NewInfluxStore(client, cfg.InfluxDBOrg, cfg.InfluxDBBucket)

    // Initialize rule engine
    ruleEngine := engine.NewRuleEngine(store)

    // Initialize API
    router := gin.Default()

    // CORS configuration
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://127.0.0.1:3000"}, // Replace with your frontend URL
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    api.SetupRoutes(router, ruleEngine)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Starting server on :%s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
