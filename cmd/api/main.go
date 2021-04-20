package main

import (
	"log"
	"os"
	"strconv"

	"fetchrewards.com/points-api/internal/db"
	"fetchrewards.com/points-api/internal/services"
	"fetchrewards.com/points-api/internal/web"
)

const DefaultPort = 8090

// Run the following from the root of the project
// go cmd/api/main.go
//
// Optional: provide a port for which to run the server
// go cmd/api/main.go 8080
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
