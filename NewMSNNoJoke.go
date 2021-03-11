package main

import (
	"bufio"
	"fmt"
	"gobweb/network"
	"log"
	"net"
	"os"
	"regexp"
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
	network.RegisterHandler("private-message", displayReceivedPrivateMessage)

	network.SetLogLevel(network.NONE)

	network.OnReady = func() {
		go startScannerOfMessenger()
	}

	if ip == nil {
		go network.StartCobweb(port)
	} else {
		go network.ConnectCobweb(port, ip, localport)
	}

	fmt.Println("Starting the New MSN, No Joke !")

	select {} // Infinite wait for something (love maybe ...)
}

// Fonction permettant de taper des messages a envoyer sur le réseau.
//
func startScannerOfMessenger() {
	var sendReg = regexp.MustCompile(`^send ((?:\d{1,3}\.){3}\d{1,3}:\d{4,5}) (.*)`)

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
		} else if sendReg.MatchString(text) {
			// Exemple :
			// send 127.0.0.1:8081 bonjour je suis un multi mot
			// 0 :  send 127.0.0.1:8081 bonjour je suis un multi mot
			// 1 :  127.0.0.1:8081
			// 2 :  bonjour je suis un multi mot

			params := sendReg.FindStringSubmatch(text)
			ipport := params[1]
			message := params[2]

			connection, err := network.GetConnection(ipport)

			if err != nil {
				log.Println("Veuillez vérifier le couple ip:port que vous avez saisi, aucune entrée trouvée")
				return
			}

			connection.Send(network.CreatePacket("private-message", message))
			return
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
	fmt.Println("[", packet.GetConnection().GetIpPortAddress(), "] -> ", packet.Pdata)
}

// Méthode qui affiche le message privé reçu.
//
func displayReceivedPrivateMessage(packet network.Packet) {
	fmt.Println("(", packet.GetConnection().GetIpPortAddress(), ") -> ", packet.Pdata)
}