package main

import (
	"fetchrewards.com/points-api/internal/db"
	"fetchrewards.com/points-api/internal/services"
	"fetchrewards.com/points-api/internal/web"
	"log"
	"os"
	"strconv"
)

const DefaultPort = 8090

func main() {
	db := db.NewInMemoryDB()
	service := services.NewPointService(db)
	server := web.NewServer(service)

	server.Start(getPort())
}

func getPort() int {
	port := DefaultPort
	if len(os.Args) > 1 {
		requestedPort, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Printf("Invalid port number %s, falling back to default", os.Args[1])
		} else {
			port = requestedPort
		}
	}
	return port
}
