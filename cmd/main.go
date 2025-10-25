package main

import (
	"fmt"
	"log"
	"os"

	"net-cat/internal/server"
)

func main() {
	port := "8989" // Порт по умолчанию
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	s, err := server.NewServer(port) // Создание сервера
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	defer s.Close()

	fmt.Println("Listening on the port :" + port)
	fmt.Println("To connect to the user use the command: nc [ip][port]")

	if err := s.Run(); err != nil { // Запуск сервера
		log.Fatalf("Error running server: %v", err)
	}
}
