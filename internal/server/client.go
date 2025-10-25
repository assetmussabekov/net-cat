package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// linuxLogo - приветственное сообщение с логотипом
var linuxLogo = "Welcome to TCP-Chat!\n" +
	"         _nnnn_\n" +
	"        dGGGGMMb\n" +
	"       @p~qp~~qMb\n" +
	"       M|@||@) M|\n" +
	"       @,----.JM|\n" +
	"      JS^\\__/  qKL\n" +
	"     dZP        qKRb\n" +
	"    dZP          qKKb\n" +
	"   fZP            SMMb\n" +
	"   HZM            MMMM\n" +
	"   FqM            MMMM\n" +
	" __| \".        |\\dS\"qML\n" +
	" |    `.       | `' \\Zq\n" +
	"_)      \\.___.,|     .'\n" +
	"\\____   )MMMMMP|   .'\n" +
	"     `-'       `--'\n" +
	"[ENTER YOUR NAME]: "

// Client - структура клиента
type Client struct {
	conn net.Conn
	name string
}

// handleConnection - обработка нового подключения
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	_, err := conn.Write([]byte(linuxLogo)) // Отправка приветствия
	if err != nil {
		log.Printf("Error sending welcome message: %v", err)
		return
	}

	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n') // Чтение имени
	if err != nil {
		log.Printf("Error reading name: %v", err)
		return
	}
	name = strings.TrimSpace(name)
	if name == "" {
		conn.Write([]byte("Name cannot be empty\n")) // Проверка пустого имени
		return
	}

	client := &Client{conn: conn, name: name}
	s.mutex.Lock()
	if s.isNameTaken(name, client) { // Проверка уникальности имени
		conn.Write([]byte("This name is already taken. Please choose another one.\n"))
		s.mutex.Unlock()
		return
	}
	s.mutex.Unlock()

	s.clientChan <- client // Добавление клиента в канал

	_, err = conn.Write([]byte("You are now connected to the chat! Use /nick <newname> to change name, /private <nick> <message> for private messages.\n"))
	if err != nil {
		log.Printf("Error sending connection confirmation: %v", err)
		return
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() { // Чтение сообщений
		message := strings.TrimSpace(scanner.Text())
		if message != "" {
			// Проверка на ANSI-последовательности (начинаются с ^[ или \x1b)
			if strings.HasPrefix(message, "\x1b[") || strings.HasPrefix(message, "^[") {
				continue
			}
			if strings.HasPrefix(message, "/nick") {
				newName := strings.TrimSpace(strings.TrimPrefix(message, "/nick"))
				if newName == "" {
					conn.Write([]byte("New name cannot be empty\n")) // Проверка пустого нового имени
					continue
				}
				s.mutex.Lock()
				if s.isNameTaken(newName, client) {
					conn.Write([]byte(fmt.Sprintf("Name '%s' is already taken by another user\n", newName)))
					s.mutex.Unlock()
					continue
				}
				oldName := client.name
				client.name = newName // Смена имени
				s.mutex.Unlock()
				s.broadcastMessage(fmt.Sprintf("%s is now known as %s", oldName, newName))
				log.Printf("%s changed name to %s", oldName, newName)
			} else if strings.HasPrefix(message, "/private") {
				parts := strings.SplitN(strings.TrimSpace(strings.TrimPrefix(message, "/private")), " ", 2)
				if len(parts) < 2 {
					conn.Write([]byte("Usage: /private <nick> <message>\n")) // Проверка формата
					continue
				}
				targetName := parts[0]
				privateMsg := parts[1]
				if privateMsg == "" {
					conn.Write([]byte("Private message cannot be empty\n")) // Проверка пустого сообщения
					continue
				}
				s.sendPrivateMessage(client, targetName, privateMsg) // Отправка приватного сообщения
			} else {
				timestamp := time.Now().Format("2006-01-02 15:04:05")
				formattedMsg := fmt.Sprintf("[%s][%s]:%s", timestamp, client.name, message)
				log.Printf("Received from %s: %s", client.name, message)
				s.broadcast <- formattedMsg // Отправка публичного сообщения
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from %s: %v", client.name, err)
	}
	s.quitChan <- client // Удаление клиента
}

// isNameTaken - проверка, занято ли имя
func (s *Server) isNameTaken(name string, currentClient *Client) bool {
	for client := range s.clients {
		if client != currentClient && strings.EqualFold(client.name, name) {
			return true
		}
	}
	return false
}

// sendPrivateMessage - отправка приватного сообщения
func (s *Server) sendPrivateMessage(sender *Client, targetName, message string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedMsg := fmt.Sprintf("[%s][Private from %s]:%s", timestamp, sender.name, message)
	senderMsg := fmt.Sprintf("[%s][Private to %s]:%s", timestamp, targetName, message)

	var target *Client
	for client := range s.clients { // Поиск получателя
		if strings.EqualFold(client.name, targetName) {
			target = client
			break
		}
	}

	if target == nil {
		sender.conn.Write([]byte(fmt.Sprintf("User '%s' not found\n", targetName))) // Ошибка: пользователь не найден
		return
	}

	_, err := target.conn.Write([]byte(formattedMsg + "\n")) // Отправка получателю
	if err != nil {
		sender.conn.Write([]byte(fmt.Sprintf("Error sending message to '%s'\n", targetName)))
		log.Printf("Error sending private message to %s: %v", targetName, err)
		return
	}

	_, err = sender.conn.Write([]byte(senderMsg + "\n")) // Подтверждение отправителю
	if err != nil {
		log.Printf("Error sending private message confirmation to %s: %v", sender.name, err)
		return
	}

	s.logToFile(formattedMsg) // Логирование
	log.Printf("Private message from %s to %s: %s", sender.name, targetName, message)
}
