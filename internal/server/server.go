package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// Константы для сервера
const (
	maxClients = 10
	logFile    = "chat_log.log"
)

// Server - структура сервера
type Server struct {
	clients    map[*Client]bool
	messages   []string
	broadcast  chan string
	mutex      sync.Mutex
	clientChan chan *Client
	quitChan   chan *Client
	logFile    *os.File
	listener   net.Listener
}

// NewServer - создание нового сервера
func NewServer(port string) (*Server, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %v", err) // Ошибка открытия лог-файла
	}

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("error starting server: %v", err) // Ошибка запуска сервера
	}

	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan string, 100),
		clientChan: make(chan *Client),
		quitChan:   make(chan *Client),
		logFile:    file,
		listener:   listener,
	}, nil
}

// Close - закрытие сервера
func (s *Server) Close() {
	s.listener.Close()
	s.logFile.Close()
}

// Run - запуск сервера
func (s *Server) Run() error {
	go s.run() // Запуск основного цикла в горутине
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handleConnection(conn) // Обработка нового клиента
	}
}

// run - основной цикл обработки событий
func (s *Server) run() {
	for {
		select {
		case client := <-s.clientChan: // Новый клиент
			s.mutex.Lock()
			if len(s.clients) >= maxClients {
				client.conn.Write([]byte("Server is full. Try again later.\n"))
				client.conn.Close() // Сервер полон
			} else {
				s.clients[client] = true
				s.broadcastMessage(fmt.Sprintf("%s has joined our chat...", client.name))
				s.sendHistory(client)
				log.Printf("Client %s joined. Total clients: %d", client.name, len(s.clients))
			}
			s.mutex.Unlock()

		case client := <-s.quitChan: // Клиент вышел
			s.mutex.Lock()
			delete(s.clients, client)
			s.broadcastMessage(fmt.Sprintf("%s has left our chat...", client.name))
			log.Printf("Client %s left. Total clients: %d", client.name, len(s.clients))
			client.conn.Close()
			s.mutex.Unlock()

		case message := <-s.broadcast: // Широковещательное сообщение
			s.mutex.Lock()
			s.messages = append(s.messages, message)
			s.logToFile(message)
			log.Printf("Broadcasting message: %s", message)
			for client := range s.clients {
				_, err := client.conn.Write([]byte(message + "\n"))
				if err != nil {
					log.Printf("Error sending to %s: %v", client.name, err)
					s.quitChan <- client // Ошибка отправки
				} else {
					log.Printf("Sent to %s: %s", client.name, message)
				}
			}
			s.mutex.Unlock()
		}
	}
}

// broadcastMessage - отправка сообщения всем клиентам
func (s *Server) broadcastMessage(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedMsg := fmt.Sprintf("[%s]%s", timestamp, message)
	s.broadcast <- formattedMsg
}

// sendHistory - отправка истории чата клиенту
func (s *Server) sendHistory(client *Client) {
	for _, msg := range s.messages {
		_, err := client.conn.Write([]byte(msg + "\n"))
		if err != nil {
			log.Printf("Error sending history to %s: %v", client.name, err)
			return
		}
	}
}

// logToFile - запись сообщения в лог-файл
func (s *Server) logToFile(message string) {
	if _, err := s.logFile.WriteString(message + "\n"); err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
	if err := s.logFile.Sync(); err != nil {
		log.Printf("Error syncing log file: %v", err)
	}
}
