package main

import (
	"context"
	"fmt"
	"github.com/IslomSobirov/tcp-with-pow/internal/client"
	"github.com/IslomSobirov/tcp-with-pow/internal/config"
)

func main() {
	fmt.Println("start client")

	// loading config from config.json
	configData, err := config.LoadConfig("config/config.json")
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}

	// init context to pass config down
	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configData)

	address := fmt.Sprintf("%s:%d", configData.ServerHost, configData.ServerPort)

	// run client
	err = client.RunClient(ctx, address)
	if err != nil {
		fmt.Println("client error:", err)
	}
}
