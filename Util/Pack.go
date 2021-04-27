package Util

import (
	"bytes"
	"encoding/binary"
)

func Pack(data []byte, action byte) ([]byte, error) {
	var err error
	writer := bytes.NewBuffer([]byte{})
	// 请求标识
	err = binary.Write(writer, binary.BigEndian, []byte{action})
	// int转字节集
	err = binary.Write(writer, binary.BigEndian, int2Bytes(len(data)))
	err = binary.Write(writer, binary.BigEndian, []byte(data))
	return writer.Bytes(), err
}

func int2Bytes(number int) []byte {
	writer := bytes.NewBuffer([]byte{})
	 _= binary.Write(writer, binary.BigEndian, int32(number))
	return writer.Bytes()
}
