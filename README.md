# TCPChat

## 🧠 Overview

**TCPChat** is a Go-based implementation of a group chat system using TCP connections, inspired by the NetCat (`nc`) utility.

The project works **only on macOS**.  
⚠️ On **Windows**, it does **not work out of the box** — additional tools (like PowerShell modules or WSL) must be installed manually.

---

## 🔍 Features

- **TCP Connection**: Supports up to 10 concurrent client connections
- **Client Naming**: Each client must provide a unique, non-empty name
- **Message Broadcasting**: All messages are sent to all connected clients
- **Message Format**: `[timestamp][client.name]:[client.message]`
- **Chat History**: New clients receive all previous messages upon joining
- **Join/Leave Notifications**: Notifies when a client joins or leaves
- **Empty Message Filtering**: Empty messages are not broadcasted
- **Robust Disconnection**: Clients stay connected if someone disconnects
- **Default Port**: Uses port `8989` if none specified

---

## Authors

• Maksat Kapan - mkapan         
• Asset Mussabekov - amussabe   
• Daniyar Shadykhanov - dshadykh   

## 🛠 Installation

Make sure you have **Go 1.23+** installed.

```bash
# Clone the repository
git clone <repository-url>
cd TCPChat

# Build the project
go build -o TCPChat main.go
```

## 🚀 Usage

Run the Server
```bash
# Start server on default port 8989
go run .

# Start server on custom port
go run . 2525

# Invalid usage will show:
go run . 2525 localhost
[USAGE]: ./TCPChat $port
```

Connect as a Client
You can use nc (NetCat):
```bash
nc <host-ip> <port>
```

Example output:
```bash
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]:
```

## 💬 Example Interaction

Server (Port 2525):
```bash
go run . 2525
Listening on the port :2525
```
Client 1 (Yenlik):
```bash
nc localhost 2525
[ENTER YOUR NAME]: Yenlik
[2020-01-20 16:03:43][Yenlik]:hello
[2020-01-20 16:03:46][Yenlik]:How are you?
Lee has joined our chat...
[2020-01-20 16:04:32][Lee]:Hi everyone!
[2020-01-20 16:04:35][Yenlik]:great, and you?
[2020-01-20 16:04:44][Lee]:good!
[2020-01-20 16:04:50][Yenlik]:bye-bye!
Lee has left our chat...
```
Client 2 (Lee):
```bash
nc localhost 2525
[ENTER YOUR NAME]: Lee
...
```

## 📚 Technical Requirements

• Language: Go

• Concurrency: Uses goroutines, channels, and mutexes

• Packages: log, os, fmt, net, sync, time, bufio, strings

## 💡 Development Guidelines

• Follow Go good practices

• Handle errors on both server and client sides

### ✅ Bonus Features (Optional)

- 💾 **Save Chat Logs**: All chat messages are saved to a log file for future reference  
- 🔒 **Private Messaging**: Clients can send private messages using the `/private <nickname> <message>` command  
  - Example: `/private Maksat Hello, how are you?`
- ✏️ **Nickname Change**: Clients can change their nickname using the `/nick <newname>` command  
  - Example: `/nick DarkKnight`

## 🎯 Learning Outcomes

This project demonstrates:

 • Go struct manipulation

 • NetCat-like functionality

 • TCP protocols and sockets

 • Go concurrency tools (goroutines, channels, mutexes)

 • IP and port management

## 🧪 Testing Checklist

• Start server on default port

• Test invalid usage

• Start server on custom port

• Connect multiple clients

• Test message broadcasting

• Verify join/leave notifications

• Verify chat history for new clients

• Test disconnection handling
