package main

import (
	"context"
	"fmt"
	"github.com/IslomSobirov/tcp-with-pow/internal/config"
	"github.com/IslomSobirov/tcp-with-pow/internal/pkg/cache"
	"github.com/IslomSobirov/tcp-with-pow/internal/pkg/clock"
	"github.com/IslomSobirov/tcp-with-pow/internal/server"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("start server")

	// loading config
	configInst, err := config.LoadConfig("config/config.json")
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}

	// initializing context and pass data
	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configInst)
	ctx = context.WithValue(ctx, "clock", clock.SystemClock{})

	cacheInst, err := cache.InitRedisCache(ctx, configInst.CacheHost, configInst.CachePort)
	if err != nil {
		fmt.Println("error init cache:", err)
		return
	}
	ctx = context.WithValue(ctx, "cache", cacheInst)

	// seed random generator to randomize order of quotes
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// run server
	serverAddress := fmt.Sprintf("%s:%d", configInst.ServerHost, configInst.ServerPort)
	err = server.RunServer(ctx, serverAddress)
	if err != nil {
		fmt.Println("server error:", err)
	}
}
