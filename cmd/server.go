package main

import (
	"fmt"
	"log"
	"time"

	// This will automatically load/inject environment variables from a .env file
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkxro/squid/internal/cache"
	"github.com/pkxro/squid/internal/cache/redis"
	"github.com/pkxro/squid/internal/common"
	"github.com/pkxro/squid/internal/config"
	"github.com/pkxro/squid/internal/controller"
	"github.com/pkxro/squid/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Set secret config
	scfg, err := config.NewSecretConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Set operations config
	ocfg, err := config.NewOperationConfig()
	if err != nil {
		log.Fatal(err)
	}

	// New Prod Logger
	zapLogger, err := common.InitZapLogger()
	if err != nil {
		log.Fatal(err)
	}

	// Flush stdout
	defer zapLogger.Sync()

	// Create new instance of redis cache
	redis, err := redis.NewRedisCache(scfg.CacheUri, scfg.CachePassword, time.Minute*15)
	if err != nil {
		log.Fatal(err)
	}

	// Create new instance of cache manager
	cacher := cache.NewCacheManager(redis)

	// Create a new instance of a controller manager
	cm := controller.NewControllerManager(
		zapLogger,
		scfg,
		ocfg,
		cacher,
	)

	// Create a new instance of a router manager
	rm := router.NewRouterManager(
		cm,
	)

	// Initialize the router
	err = rm.InitRouter()
	if err != nil {
		log.Fatal(err)
	}

	zapLogger.Info("Starting Squid...")

	// Start the server
	err = rm.Router.Run(fmt.Sprintf("0.0.0.0%v", common.FormatPort(*&ocfg.Port)))
	if err != nil {
		log.Fatal(err)
	}
}
