package main

import (
	"context"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/eine-doodka/twoStage/cache"
	"github.com/eine-doodka/twoStage/config"
	"github.com/eine-doodka/twoStage/logic"
	"github.com/eine-doodka/twoStage/server"
	"github.com/go-redis/redis/v9"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/twoStage.toml", "cfg file path")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()
	cfg := config.Default()
	_, err := toml.DecodeFile(configPath, &cfg)
	if err != nil {
		log.Fatalln("Config read error:", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatalln("Redis init error:", err)
	}
	cacheModule := cache.NewImpl(
		redisClient,
		cfg.StorageDuration,
	)
	logicModule := logic.NewImpl(
		cacheModule,
		cfg.CodeLength,
	)
	handlersModule := server.NewHandlers(logicModule)
	routesModule := server.NewServer(handlersModule)
	err = http.ListenAndServe(":8080", routesModule)
	if err != nil {
		log.Fatalln("Server start error:", err)
	}
}
