package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zekeriyyah/lujay-autocity/internal/config"
	"github.com/zekeriyyah/lujay-autocity/internal/routes"
	"github.com/zekeriyyah/lujay-autocity/pkg"
)


func main() {
	r := gin.Default()

	cfg, err := config.LoadConfig()
	if err != nil {
		pkg.Error(err, "failed to load configurations")
		return
	}

	r = routes.SetupRouter(r, cfg)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000" 
	}

	pkg.Info("Server starting on port " + port)

	if err := r.Run(":" + port); err != nil {
		pkg.Error(err, "failed to start server")
		return
	}
}
