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

	network.OnReady = func() {
		go startScannerOfMessenger()
	}

	if ip == nil {
		go network.StartCobweb(port)
	} else {
		go network.ConnectCobweb(port, ip, localport)
	}

	fmt.Println("Starting the New MSN, No Joke !")
	fmt.Println("[IP:PORT] -> MESSAGE is for broadcast messages")
	fmt.Println("(IP:PORT) -> MESSAGE is for private messages")

	select {} // Infinite wait for something (love maybe ...)
}

// Fonction permettant de taper des messages a envoyer sur le réseau.
//
func startScannerOfMessenger() {
	var sendReg = regexp.MustCompile(`^send ((?:\d{1,3}\.){3}\d{1,3}:\d{4,5}) (.*)`)
	var loglevelReg = regexp.MustCompile(`^set loglevel (INFO|WARNING|ERROR|NONE)`)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')

		if err != nil {
			continue
		}

		if text == "commands\n" {
			fmt.Println("Help :")
			fmt.Println("\tcommands: Displays this message")
			fmt.Println("\tlist: Lists all connections")
			fmt.Println("\tsend <IP:PORT> <MESSAGE>: Sends a private message to IP:PORT connection")
			fmt.Println("\tset loglevel <INFO|WARNING|ERROR|NONE>: Sets the loglevel (default: INFO)")
			continue
		} else if text == "list\n" {
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
				continue
			}

			connection.Send(network.CreatePacket("private-message", message))
			continue
		} else if loglevelReg.MatchString(text) {
			params := loglevelReg.FindStringSubmatch(text)

			loglevel := params[1]

			var loglevelInt int
			switch loglevel {
			case "INFO":
				loglevelInt = network.INFO
				break
			case "WARNING":
				loglevelInt = network.WARNING
				break
			case "ERROR":
				loglevelInt = network.ERROR
				break
			case "NONE":
				loglevelInt = network.NONE
				break
			default:
				return
			}
			fmt.Printf("Log level set to %s\n", loglevel)
			network.SetLogLevel(loglevelInt)
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
	fmt.Println("[", packet.GetConnection().GetIpPortAddress(), "] -> ", packet.Pdata)
}

// Méthode qui affiche le message privé reçu.
//
func displayReceivedPrivateMessage(packet network.Packet) {
	fmt.Println("(", packet.GetConnection().GetIpPortAddress(), ") -> ", packet.Pdata)
}