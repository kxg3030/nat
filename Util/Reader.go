package Util

import (
	"encoding/binary"
	"errors"
	"nat/Logger"
	"net"
)

type Reader struct {
	Buff       []byte
	Start      int
	End        int
	BuffLen    int
	HeaderLen  int
	BodyOffset int
	Content    chan []byte
	Connector  net.Conn
}

func (r *Reader) Read() error {
	for {
		r.moveOffset()
		if r.End >= r.BuffLen {
			Logger.Logger.Println("data is too long")
			close(r.Content)
			return errors.New("data is too long")
		}
		// 读取数据到缓冲区(1kb)
		length, err := r.Connector.Read(r.Buff[r.End:])
		if err != nil {
			Logger.Logger.Println("read data error：" + err.Error())
			close(r.Content)
			return err
		}
		// 缓冲区下次填充的位置后移
		r.End += length
		// 切割完整数据
		r.readOneMessage()
	}
}
func (r *Reader) readOneMessage() {
	// 检查缓冲区数据是否完整
	if r.End-r.Start < r.HeaderLen {
		// 头部不足
		return
	}
	// 读取包头部
	headerData := r.Buff[r.Start : r.HeaderLen+r.Start]
	// 读取包头中包体长度
	bodyLength := binary.BigEndian.Uint16(headerData[r.BodyOffset : r.BodyOffset+2])
	// 判断包体的长度
	if r.End-r.Start-r.HeaderLen < int(bodyLength) {
		// 包体不足
		return
	}
	// 完整的包头包体
	r.Content <- r.Buff[r.Start : r.Start+r.HeaderLen+int(bodyLength)]
	// 下一次读取的开始位置
	r.Start += r.HeaderLen + int(bodyLength)
	r.readOneMessage()
}

func (r *Reader) moveOffset() {
	if r.Start == 0 {
		return
	}
	// 将缓冲区的不完整数据保存在buffer中
	copy(r.Buff, r.Buff[r.Start:r.End])
	// 这里计算新缓冲区的填充位置(其实就是计算剩余的那一部分数据的长度是多少)
	r.End -= r.Start
	r.Start = 0
}