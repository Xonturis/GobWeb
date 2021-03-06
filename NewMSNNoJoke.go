package main

import (
	"bufio"
	"fmt"
	"gobweb/network"
	"net"
	"os"
)

//func Test() {
//	time.Sleep(2 * time.Second)
//	network.SendToAll(network.CreatePacket("127.0.0.1", []byte("BONJOUR!")))
//	network.SendToAll(network.CreatePacket("127.0.0.1", []byte("YES!")))
//}

func LaunchNMSNNJ(sport int, ip net.IP, localport int) {

	network.RegisterHandler("message", displayReceivedMessage)

	network.OnReady = func() {
		go startScannerOfMessenger()
	}

	if ip == nil {
		go startServer(sport)
	} else {
		go connectToCobweb(sport, ip, localport)
	}

	select {}
}

func startScannerOfMessenger() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		network.SendToAll(network.CreatePacket("message", text))
	}
}

func displayReceivedMessage(packet network.Packet) {
	fmt.Println("[", packet.PipSrc, "] ", packet.Pdata)
}

func startServer(sport int) {
	network.StartCobweb(sport)
}

func connectToCobweb(sport int, ip net.IP, localport int) {
	network.ConnectCobweb(sport, ip, localport)
}
