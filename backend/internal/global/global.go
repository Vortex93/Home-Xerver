package global

import (
	"os"

	"github.com/joho/godotenv"
)

// ========================================
// State
// ========================================
var (
    API_HOST string
    API_PORT string
)


// ========================================
// Bootstrap
// ========================================
func init() {
    err := godotenv.Load()
    if err != nil {
        panic("Error loading .env file: " + err.Error())
    }

    API_HOST = os.Getenv("API_HOST")
    API_PORT = os.Getenv("API_PORT")
}

// ========================================
// Functions
// ========================================

// GetAPIHost returns the API host from environment variables.
func GetAPIHost() string {
    if API_HOST == "" {
        return "localhost" // Default to localhost if not set
    }
    return API_HOST
}

// GetAPIPort returns the API port from environment variables.
func GetAPIPort() string {
    if API_PORT == "" {
        return "8080" // Default to 8080 if not set
    }
    return API_PORT
}

// ========================================
// Types
// ========================================
