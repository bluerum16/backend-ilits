package main

import (
	"log"
	"os"

	"backend-ilits/config"
	"backend-ilits/routes"
)

func main() {
	db := config.ConnectDatabase()
	r := routes.SetupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server jalan di port " + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
