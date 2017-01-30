package ddb

import (
	"net"
	"log"
	"os"
	"encoding/json"
	"strings"
)

var stdlog, errlog *log.Logger

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

// Accept a client connection and collect it in a channel
func AcceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		stdlog.Println("Client ", conn.RemoteAddr(), " just established Connection")
		listen <- conn
	}
}

func HandleClient(client net.Conn) {
	var cmd string
	for {
		buf := make([]byte, 4096)
		numBytes, err := client.Read(buf)

		if numBytes == 0 || err != nil {
			client.Close()
			return
		}

		cleanedUpCommands := strings.TrimSpace(string(buf))
		commands := strings.Split(cleanedUpCommands, "\r\n")
		for _, cmd = range commands {

			if !isEmptyCommand(cmd) {
				handleCommand(client, cmd)
			}

		}
	}
}

func isEmptyCommand(command string) bool {

	c := []byte(command)
	for _, val := range c {
		if val != 0 {
			return false
		}
	}

	return true
}

func handleCommand(client net.Conn, command string) {
	cmd := parseCommand(command)

	switch cmd.operation {
	case "read":
		db := getInstance()
		target, ok := db.targets[cmd.target]
		var data map[string]interface{}

		if ok == false {
			data = make(map[string]interface{})
		} else {
			data = target.read(cmd.data);
		}

		stdlog.Println("Data readed from db: ", data)
		jsonData, err := json.Marshal(data)
		if err == nil {
			client.Write(jsonData)
		} else {
			client.Write([]byte("[]"))
		}
		client.Write([]byte("\n"))

	case "write":
		db := getInstance()
		target, ok := db.targets[cmd.target]
		if ok == false {
			target = Target{name: cmd.target, data: make(map[string]interface{})}
			db.targets[cmd.target] = target
		}
		stdlog.Println("Data to be inserted: ", cmd.data)
		id := target.insert(cmd.data)
		stdlog.Println("Data inserted with id", id)

		data := make(map[string]string);
		data["_id"] = id;
		jsonData, err := json.Marshal(data)
		if err == nil {
			client.Write(jsonData)
		} else {
			client.Write([]byte("[]"))
		}
		client.Write([]byte("\n"))
	case "exit":
		stdlog.Println("Client ", client.RemoteAddr(), " exited Connection")
		client.Close()
	default:
		errlog.Println("Unknown Operation received: '", cmd.operation, "'")
		client.Close()
	}
}

func parseCommand(command string) Command {

	var c Command

	var f interface{}
	err := json.Unmarshal([]byte(command), &f)

	if err != nil {
		errlog.Println("Unknown command received: ", command, "[", err, "]")
	}

	c.initFromJson(f)

	return c
}
