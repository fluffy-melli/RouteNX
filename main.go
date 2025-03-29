package main

import (
	"fmt"
	"log"
	"net"

	"github.com/fluffy-melli/RouteNX/internal/middleware"
	"github.com/fluffy-melli/RouteNX/pkg/config"
)

func main() {
	configs, err := config.LoadFromFile(config.RouteNXJSON)
	if err != nil {
		configs = config.NewRouteNX()
		configs.SaveToFile(config.RouteNXJSON)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", configs.Port))
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go middleware.Listener(configs, conn)
	}
}
