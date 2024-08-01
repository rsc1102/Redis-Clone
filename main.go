package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	port := "6379"
	fmt.Println("Client listening at port:", port)

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := aof.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		if command == "SHUTDOWN" { // command to shutdown the redis server
			fmt.Println("Shutting down redis server...")
			break
		}

		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		if command == "SET" || command == "HSET" || command == "DEL" || command == "HDEL" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)

	}
}
