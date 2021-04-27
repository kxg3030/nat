package Server

import (
	"encoding/json"
	"fmt"
	"nat/Constant"
	"nat/Logger"
	"nat/Util"
	"net"
	"strconv"
	"time"
)

// 请求标识
const ReqLength = 1

// 包体长度字节
const BodyLength = 4

// 包头长度
const HeaderLength = ReqLength + BodyLength

// 包头标识包体长度的起始位置
const BodyOffset = 1

type Server struct {
	Nodes map[string]*Node
}

type Node struct {
	createdAt    int64
	lastActiveAt int64
	connector    net.Conn
	reader       *Util.Reader
}

func NewServer() *Server {
	return &Server{
		Nodes: make(map[string]*Node),
	}
}

// 外部服务监听端口
func (s *Server) server4Net() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:9501")
	if err != nil {
		Logger.Logger.Println("resolve address error：" + err.Error())
		return
	}
	Logger.Logger.Println("server4Net start success")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		Logger.Logger.Println("server4Net listen error：" + err.Error())
		return
	}
	for {
		connect, err := listener.Accept()
		if err != nil {
			Logger.Logger.Println("server4Net accept error：" + err.Error())
			return
		}
		// 客户端
		connector := &Node{
			createdAt:    time.Now().Unix(),
			lastActiveAt: time.Now().Unix(),
			connector:    connect,
			reader: &Util.Reader{
				Start:      0,
				End:        0,
				BuffLen:    1024,
				HeaderLen:  HeaderLength,
				BodyOffset: BodyOffset,
				Content:    make(chan []byte, 10),
				Connector:  connect,
			},
		}
		Logger.Logger.Println("new connect：" + connect.RemoteAddr().String() + " nodes：" + strconv.Itoa(len(s.Nodes)))
		s.Nodes[connect.RemoteAddr().String()] = connector
		go s.handleNode(connector)
	}
}

// 客户端服务监听端口
func (s *Server) server4Local() {

}

// 监听外部服务
func (s *Server) Lister() {
	s.server4Net()
}

// 处理客户端
func (s *Server) handleNode(node *Node) {
	defer func() {
		node.connector.Close()
		key := node.connector.RemoteAddr().String()
		if _, ok := s.Nodes[key]; ok {
			delete(s.Nodes, key)
		}
	}()
	// 读取数据
	node.reader.Buff = make([]byte, node.reader.BuffLen)
	go node.reader.Read()
	// 读取
	for {
		select {
		case data, ok := <-node.reader.Content:
			if !ok {
				fmt.Println("client quit")
				return
			}
			node.service(data)
		}
	}
}

func (n *Node) service(data []byte) {
	action := data[0]
	body := data[HeaderLength:]
	var bodyMap map[string]interface{}
	if len(body) > 0 {
		err := json.Unmarshal(body, &bodyMap)
		if err != nil {
			Logger.Logger.Println("json unmarshal error：" + err.Error())
			return
		}
	}
	switch action {
	case Constant.Login:
		// 验证登录
		username := bodyMap["username"].(string)
		password := bodyMap["password"].(string)
		fmt.Println(username, password)
		_, _ = n.connector.Write([]byte("login success"))
	case Constant.NodeList:
		_, _ = n.connector.Write([]byte("select success"))
	}
}
