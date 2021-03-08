package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

//Démarrage de l'application et du réseau.
func main() {
	args := os.Args[1:]
	lenArgs := len(args)

	// Juste un port en param
	if lenArgs >= 1 {
		//get le port
		port, err := strconv.Atoi(args[0])
		if err == nil {
			if lenArgs >= 2 {
				ipPort := strings.Split(args[1], ":")
				if len(ipPort) == 2 {
					ip := net.ParseIP(ipPort[0])
					port2, err := strconv.Atoi(ipPort[1])
					if ip != nil && err == nil {
						LaunchNMSNNJ(ip, port2, port)
					} else {
						fmt.Println("Erreur, l'adresse IP n'est pas valide")
						fmt.Println("Syntaxe :")
						fmt.Println("cobweb <port> [ip:port]")
						os.Exit(3)
					}
				} else {
					fmt.Println("Erreur, le format ip:port n'est pas respecté")
					fmt.Println("Syntaxe :")
					fmt.Println("cobweb <port> [ip:port]")
					os.Exit(3)
				}
			} else {
				LaunchNMSNNJ(nil, port, 0)
			}

		} else {
			fmt.Println("Erreur, le port n'est pas valide")
			fmt.Println("Syntaxe :")
			fmt.Println("cobweb <port> [ip:port]")
			os.Exit(2)
		}

		// Erreur, pas assez de params
	} else {
		fmt.Println("Syntaxe :")
		fmt.Println("cobweb <port> [ip:port]")
		os.Exit(1)
	}
}
