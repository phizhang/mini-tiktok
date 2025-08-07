package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phizhang/mini-tiktok/internal/api"
	"github.com/phizhang/mini-tiktok/internal/configs"
	"github.com/phizhang/mini-tiktok/internal/db"
)

func main() {
	configs.LoadEnv()
	log.Printf("Loaded environment: %+v", configs.GlobalEnv)

	// Initialize Cassandra database connection
	dbHost := configs.GlobalEnv.DBHost
	if dbHost == "" {
		log.Fatal("DB_HOSTS environment variable is required (comma separated list)")
	}
	hosts := strings.Split(dbHost, ",")
	keyspace := "mini_tiktok"
	if err := db.InitCassandra(hosts, keyspace); err != nil {
		log.Fatalf("Failed to initialize Cassandra: %v", err)
	}
	defer db.CloseCassandra()

	router := gin.Default()

	router.POST("/api/upload", api.UploadVideoHandler)

	log.Println("Starting server on :8080...")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
