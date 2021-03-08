package main

import (
	"bufio"
	"fmt"
	"gobweb/network"
	"net"
	"os"
)

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
		text, err := reader.ReadString('\n')

		if err != nil {
			continue
		}

		if text == "list\n" {
			listAllConnections()
			continue
		}
		network.SendToAll(network.CreatePacket("message", text))
	}
}

func listAllConnections() {
	for _, ipport := range network.GetAllConnectedIPListeningPortString() {
		fmt.Println(ipport)
	}
}

func displayReceivedMessage(packet network.Packet) {
	fmt.Println("[", packet.Conn.LocalAddr, "] ", packet.Pdata)
}

func startServer(sport int) {
	network.StartCobweb(sport)
}

func connectToCobweb(sport int, ip net.IP, localport int) {
	network.ConnectCobweb(sport, ip, localport)
}
