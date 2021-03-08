package main

import (
	"bufio"
	"fmt"
	"gobweb/network"
	"net"
	"os"
)

//
// Ce programme nous a servi de programme d'essai pour débuger et tester en situation réelle le réseau.
//

// Méthode permettant de lancer l'application NewMSNNoJoke
//
//  net.IP  ip         L'ip sur laquelle se connecter.
//  int     port       Le port de l'ip sur laquelle se connecter.
//  int     localport  Le port local d'écoute.
//
func LaunchNMSNNJ(ip net.IP, port int, localport int) {

	network.RegisterHandler("message", displayReceivedMessage)

	network.OnReady = func() {
		go startScannerOfMessenger()
	}

	if ip == nil {
		go network.StartCobweb(port)
	} else {
		go network.ConnectCobweb(port, ip, localport)
	}

	select {} // Infinite wait for something (love maybe ...)
}

// Fonction permettant de taper des messages a envoyer sur le réseau.
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

// Méthode écrivant sur la console toutes les connexions du pair.
//
func listAllConnections() {
	for _, ipport := range network.GetAllConnectedIPListeningPortString() {
		fmt.Println(ipport)
	}
}

// Méthode qui affiche le message reçu.
//
func displayReceivedMessage(packet network.Packet) {
	fmt.Println("[", packet.GetConnection().GetConn().RemoteAddr().String(), "] ", packet.Pdata)
}
