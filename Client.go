package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"nat/Logger"
	"net"
)

func main() {
	data, _ := bag()
	connect, err := net.Dial("tcp", "127.0.0.1:9501")
	if err != nil {
		Logger.Logger.Println(err.Error())
		return
	}
	Logger.Logger.Println("connect server success")
	for i := 0; i < 50; i++ {
		buff := make([]byte, 1024)
		_, _ = connect.Write(data)
		_, err = connect.Read(buff[0:])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(buff))
	}
	connect.Close()
}

func bag() ([]byte, error) {
	var err error
	text, _ := json.Marshal(map[string]string{
		"username": "kxg3030",
		"password": "xxx",
	})
	writer := bytes.NewBuffer([]byte{})
	// 请求标识
	err = binary.Write(writer, binary.BigEndian, []byte{0x01})
	// int转字节集
	err = binary.Write(writer, binary.BigEndian, int2Bytes(len(text)))
	err = binary.Write(writer, binary.BigEndian, []byte(text))
	return writer.Bytes(), err
}

func int2Bytes(number int) []byte {
	writer := bytes.NewBuffer([]byte{})
	_ = binary.Write(writer, binary.BigEndian, &number)
	return writer.Bytes()
}
