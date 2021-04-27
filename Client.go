package main

import (
	"bytes"
	"encoding/binary"
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
	text := "1主机字节序模式有两种，大端数据模式和小端数据模式，在网络编程中应注意这两者的区别，以保证数据处理的正确性；例如网络的数据是以大端数据模式进行交互，而我们的主机大多数以小端模式处理，如果不转换，数据会混乱 参考 ；一般来说，两个主机在网络通信需要经过如下转换过程：主机字节序 —> 网络字节序 -> 主机字节序0"
	writer := bytes.NewBuffer([]byte{})
	err = binary.Write(writer, binary.BigEndian, []byte{0x01})
	// int转字节集
	bytesBuffer := bytes.NewBuffer([]byte{})
	err = binary.Write(bytesBuffer, binary.BigEndian, int32(len(text)))
	err = binary.Write(writer, binary.BigEndian, bytesBuffer.Bytes())
	err = binary.Write(writer, binary.BigEndian, []byte(text))
	fmt.Println(writer.Bytes())
	return writer.Bytes(), err
}
