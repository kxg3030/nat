package main

import (
	"encoding/json"
	"fmt"
	"nat/Logger"
	"nat/Util"
	"net"
)

func main() {
	text, _ := json.Marshal(map[string]string{
		"username": "kxg3030",
		"password": "xxx",
	})
	data, _ := Util.Pack(text,0x03)
	connect, err := net.Dial("tcp", "127.0.0.1:9502")
	if err != nil {
		Logger.Logger.Println(err.Error())
		return
	}
	Logger.Logger.Println("connect server success")
	buff := make([]byte, 1024)
	_, _ = connect.Write(data)
	_, err = connect.Read(buff[0:])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buff))
	connect.Close()
}
