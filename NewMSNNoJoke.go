package main

import (
	"fmt"
	"gobweb/network"
	"net"
)

//func Test() {
//	time.Sleep(2 * time.Second)
//	fmt.Println("SENDING")
//	network.SendToAll(network.CreatePacket("127.0.0.1", []byte("BONJOUR!")))
//	network.SendToAll(network.CreatePacket("127.0.0.1", []byte("YES!")))
//}


func LaunchNMSNNJ(sport int, ip net.IP, localport int) {

	network.RegisterHandler("message", func(packet network.Packet) {
		fmt.Println(packet)
	})

	network.OnReady = func() {
		go startScannerOfMessenger()
	}

	if ip == nil {
		go startServer(sport)
	} else {
		go connectToCobweb(sport, ip, localport)
	}


	for {}
}


func startScannerOfMessenger() {
	fmt.Println("qdsfqsdfqsdf")
	network.SendToAll(network.CreatePacket("127.0.0.1", "message", "BONJOUR!"))

}


func startServer(sport int) {
	network.StartCobweb(sport)
}

func connectToCobweb(sport int, ip net.IP, localport int) {
	network.ConnectCobweb(sport, ip, localport)
}
