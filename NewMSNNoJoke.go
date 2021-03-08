package main

import (
	"bufio"
	"fmt"
	"gobweb/network"
	"net"
	"os"
)

// Méthode permettant de lancer l'application NewMSNNoJoke
//
//  net.IP ip         L'ip sur laquelle se connecter.
//  int    port       Le port de l'ip sur laquelle se connecter.
//  int    localport  Le port local d'écoute.
//
func LaunchNMSNNJ(ip net.IP, port int, localport int) {

	network.RegisterHandler("message", displayReceivedMessage)

	network.OnReady = func() {
		go startScannerOfMessenger()
	}

	if ip == nil {
		go startServer(port)
	} else {
		go connectToCobweb(ip, port, localport)
	}

	select {}
}

// Méthode permettant de taper des messages a envoyer sur le réseau.
//
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

func connectToCobweb(ip net.IP, port int, localport int) {
	network.ConnectCobweb(port, ip, localport)
}
