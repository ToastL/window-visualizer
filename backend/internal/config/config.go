package config

import (
	"os"
	"time"
)

type Config struct {
	Port           string
	BroadcastRate  time.Duration
	WriteDeadline  time.Duration
	ChromeOffsetY  int
	WindowCacheTTL time.Duration
	PongWait       time.Duration
	PingPeriod     time.Duration
}

var AppConfig = Config{
	Port:           getEnv("PORT", "8080"),
	BroadcastRate:  4 * time.Millisecond,
	WriteDeadline:  200 * time.Millisecond,
	ChromeOffsetY:  90,
	WindowCacheTTL: 100 * time.Millisecond,
	PongWait:       60 * time.Second,
	PingPeriod:     54 * time.Second,
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
