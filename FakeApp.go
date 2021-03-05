package main

import (
	"cobweb/network"
	"fmt"
	"time"
)

func Test() {
	time.Sleep(2 * time.Second)
	fmt.Println("SENDING")
	network.SendToAll(network.CreatePacket("127.0.0.1", []byte("BONJOUR!")))
	network.SendToAll(network.CreatePacket("127.0.0.1", []byte("YES!")))
}
