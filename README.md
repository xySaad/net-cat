# Netcat Chat

Netcat Chat is a simple TCP chat application written in Go. It allows multiple users to connect to a server and communicate with each other in real-time. The application is designed to be easy to use.

## Features

- TCP-based chat functionality
- Anonymity, no authentication
- Chat history logging
- Simple command-line interface
- Handles multiple groups

## Prerequisites
- Go 1.22.3 or later for the host or clients with the custom ui
- nc command for clients

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/xySaad/net-cat.git
   cd net-cat
   ```

2. **Build the application:**

   ```bash
   go build -o TCPChat .
   ```

3. **Run the application:**

   You can specify a port number when starting the server. If no port is specified, it defaults to `8989`.

   ```bash
   ./TCPChat [port]
   ```

## Usage

1. **Start the server:**

   Run the server using the command above. For example, to start the server on port `8989`, simply run:

   ```bash
   ./TCPChat
   ```

2. **Connect to the chat:**

   Users can connect to the chat using a TCP client (like `netcat`). For example, using `netcat`:

   ```bash
   nc localhost 8989
   ```

3. **Enter your username:**

   When prompted, enter a unique username. The username must be between 3 and 25 characters long and can only contain Latin letters and hyphens.

4. **Start chatting:**

   Once connected, you can start sending messages. All users in the chat room will receive your messages.

5. **Disconnect:**

   When you want to leave the chat, simply close the connection (e.g., by closing the terminal or using `Ctrl+C`).

## Logging

The application logs chat history and connection events to the following files:

- `history/somegroup_<port>.chat.log`: Contains the chat messages for the group named 'somegroup'.
- `logs/netcat-connection_<port>.log`: Contains user connection and disconnection events and server logs.


## Contributing

Contributions are welcome! If you have suggestions for improvements or new features, feel free to open an issue or submit a pull request.

### Contributors

- [xySaad](https://github.com/xySaad)
- [0Emperor](https://github.com/0Emperor)