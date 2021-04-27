package Server

import (
	"encoding/json"
	"fmt"
	"io"
	"nat/Constant"
	"nat/Logger"
	"nat/Util"
	"net"
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

const Server4Net = "127.0.0.1:9502"

const Server4Local = "127.0.0.1:9501"

type Server struct {
	Node *Node
}

type Node struct {
	createdAt    int64
	lastActiveAt int64
	connector    net.Conn
	reader       *Util.Reader
}

func NewServer() *Server {
	return &Server{

	}
}

// 外部服务监听端口
func (s *Server) server4Local() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", Server4Local)
	if err != nil {
		Logger.Logger.Println("resolve address error：" + err.Error())
		return
	}
	Logger.Logger.Println("server4Local start success")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		Logger.Logger.Println("server4Local listen error：" + err.Error())
		return
	}
	for {
		connect, err := listener.Accept()
		if err != nil {
			Logger.Logger.Println("server4Local accept error：" + err.Error())
			return
		}
		// 客户端
		Logger.Logger.Println("new connect：" + connect.RemoteAddr().String())
		s.Node = &Node{
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
		go s.handleLocalNode(s.Node)
	}
}

// 客户端服务监听端口
func (s *Server) server4Net() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", Server4Net)
	if err != nil {
		Logger.Logger.Println("resolve address error：" + err.Error())
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		Logger.Logger.Println("server4Net listen error：" + err.Error())
		return
	}
	Logger.Logger.Println("server4Net start success")
	for {
		connect, err := listener.Accept()
		if err != nil {
			Logger.Logger.Println("server4Net accept error：" + err.Error())
			return
		}
		go s.handleNetNode(connect)
	}
}

// 监听外部服务
func (s *Server) Lister() {
	go s.server4Net()
	time.Sleep(time.Second * 2)
	go s.server4Local()
}

// 处理客户端
func (s *Server) handleLocalNode(node *Node) {
	defer func() {
		node.connector.Close()
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

func (s *Server) handleNetNode(client net.Conn) {
	defer client.Close()
	buff := make([]byte, 1024)
	connect, err := net.Dial("tcp", Server4Local)
	if err != nil {
		Logger.Logger.Println(err.Error())
		return
	}
	// 转接到本地9501端口
	_, err = io.Copy(connect, client)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 读取返回
	result, _ := connect.Read(buff[0:])
	fmt.Println("back：" + string(result))
	connect.Close()
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
	case Constant.Transport:
		// 转发到局域网的连接
	}
}
