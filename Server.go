package main

import (
	"fmt"
	"nat/Server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	s := Server.NewServer()
	s.Lister()
	<-signalChan
	fmt.Println("server quit")
}
