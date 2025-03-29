package main

import (
	"fmt"
	"log"
	"net"

	"github.com/fluffy-melli/RouteNX/internal/cache"
	"github.com/fluffy-melli/RouteNX/internal/middleware"
	"github.com/fluffy-melli/RouteNX/internal/router"
	"github.com/fluffy-melli/RouteNX/pkg/config"
)

func main() {
	var err error
	cache.Config, err = config.LoadFromFile(config.RouteNXJSON)
	if err != nil {
		cache.Config = config.NewRouteNX()
		cache.Config.SaveToFile(config.RouteNXJSON)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cache.Config.Port))
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	go func() {
		router := router.Router()
		router.Run(fmt.Sprintf(":%d", cache.Config.WebPort))
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go middleware.Listener(conn)
	}
}
