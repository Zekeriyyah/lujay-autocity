package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/zekeriyyah/lujay-autocity/pkg"
)

type Config struct {
    DatabaseURL string
    JWTSecret   string
    // Add other config variables as needed
}

func LoadConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		pkg.Info("Warning: .env file could not be loaded. Relying on OS environment variables.")
	}

    return &Config{
        DatabaseURL: getEnv("DATABASE_URL", "postgresql://autocity_user:x0bHy2QL75QdSmmXBCz0h6IaCEKFZYS1@dpg-d47dg64hg0os73fju09g-a.oregon-postgres.render.com/autocity_db_cecf"),
        JWTSecret:   getEnv("JWT_SECRET", "!@#$%^&*()"),
    }, nil
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}