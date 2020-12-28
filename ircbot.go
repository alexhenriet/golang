package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

// Config ..
type Config struct {
	Bot struct {
		Nickname string
		Ident    string
		Realname string
	}
	Server struct {
		Host    string
		Port    string
		Channel string
	}
}

// Message ..
type Message struct {
	From   string
	Action string
	To     string
	Text   string
}

func main() {
	config := readConfig("ircbot-config.json")
	fmt.Printf("Connecting to %v\n", config.Server.Host+":"+config.Server.Port)
	conn, err := net.Dial("tcp", config.Server.Host+":"+config.Server.Port)
	if err != nil {
		panic(err)
	}
	connect(conn, config)
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			panic(err)
		}
		handleRawMessage(conn, config, message)
	}
}

func readConfig(filename string) Config {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config := Config{}
	err = decoder.Decode(&Config)
	if err != nil {
		panic(err)
	}
	return Config
}

func connect(conn net.Conn, config Config) {
	fmt.Fprintf(conn, "USER %s 0 * :%s\n", config.Bot.Ident, config.Bot.Realname)
	fmt.Fprintf(conn, "NICK %s\n", config.Bot.Nickname)
	fmt.Fprintf(conn, "JOIN %s\n", config.Server.Channel)
}

func handleRawMessage(conn net.Conn, config Config, rawMessage string) {
	fmt.Printf("%s", rawMessage)
	if strings.HasPrefix(rawMessage, "PING :") {
		reply(conn, strings.Replace(rawMessage, "PING :", "PONG :", -1))
	}
	if strings.HasPrefix(rawMessage, ":") {
		message := parseRawMessage(rawMessage)
		if message.Action == "INVITE" {
			reply(conn, "JOIN "+message.Text)
		}
		if message.Action == "JOIN" {
			reply(conn, "PRIVMSG "+message.To+" :Hello folks!\n")
		}
		if message.Action == "KICK" && strings.HasPrefix(message.Text, config.Bot.Nickname) {
			reply(conn, "JOIN "+message.To)
		}
	}
}

func reply(conn net.Conn, rawMessage string) {
	fmt.Printf("%s", rawMessage)
	fmt.Fprintf(conn, "%s\n", rawMessage)
}

func parseRawMessage(message string) Message {
	parts := strings.Split(message, " ")
	from := strings.TrimPrefix(parts[0], ":")
	action := parts[1]
	to := strings.TrimSpace(strings.TrimPrefix(parts[2], ":"))
	text := strings.TrimSpace(strings.TrimPrefix(strings.Join(parts[3:], " "), ":"))
	/*
		fmt.Printf("from: %s\n", from)
		fmt.Printf("action: %s\n", action)
		fmt.Printf("to: %s\n", to)
		fmt.Printf("text: %s\n", text)
	*/
	return Message{from, action, to, text}
}
