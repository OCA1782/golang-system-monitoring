package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"sysmonitor/handlers"
	"sysmonitor/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	alertConfig = models.AlertConfig{
		CPUThreshold:  80,
		MemThreshold:  85,
		DiskThreshold: 90,
		GPUThreshold:  85,
		SoundEnabled:  true,
		NotifyEnabled: true,
	}
	alertMu sync.RWMutex

	configFile = "alert_config.json"
)

func main() {
	loadConfig()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	hub := handlers.NewHub()
	hub.StartBroadcast(1 * time.Second)

	r.GET("/ws", func(c *gin.Context) {
		hub.ServeWS(c.Writer, c.Request)
	})

	r.GET("/api/processes", func(c *gin.Context) {
		procs, err := handlers.GetProcesses()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, procs)
	})

	r.GET("/api/alerts/config", func(c *gin.Context) {
		alertMu.RLock()
		defer alertMu.RUnlock()
		c.JSON(http.StatusOK, alertConfig)
	})

	r.POST("/api/alerts/config", func(c *gin.Context) {
		var cfg models.AlertConfig
		if err := c.ShouldBindJSON(&cfg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		alertMu.Lock()
		alertConfig = cfg
		alertMu.Unlock()
		saveConfig(cfg)
		c.JSON(http.StatusOK, gin.H{"status": "saved"})
	})

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Unix()})
	})

	log.Println("SysMonitor backend running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return
	}
	alertMu.Lock()
	defer alertMu.Unlock()
	_ = json.Unmarshal(data, &alertConfig)
}

func saveConfig(cfg models.AlertConfig) {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	_ = os.WriteFile(configFile, data, 0644)
}
