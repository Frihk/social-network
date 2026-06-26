package main

<<<<<<< HEAD
import "fmt"

func main()  {
	fmt.Print("social")	
=======
import (
	"log"

	"social/pkg/db/sqlite"
	"social/server"
)

func main() {
	if err := sqlite.InitDB("./social-network.db", "file://pkg/db/migration/sqlite"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer sqlite.CloseDB()

	srv := server.NewServer()
	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
>>>>>>> kagunda
}