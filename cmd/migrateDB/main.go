package main

import (
	"github.com/zekeriyyah/lujay-autocity/internal/config"
	"github.com/zekeriyyah/lujay-autocity/internal/database"
)

func main () {
	envVar, _ := config.LoadConfig()
	
	database.InitDB(envVar.DatabaseURL)
}