package main

import (
	"flag"
	"net"
)

//Démarrage de l'application et du réseau.
func main() {
	port := flag.Int(
		"port", 8080,
		"Port sur lequel écouter")
	isFirst := flag.Bool(
		"first", false,
		"Pour déclencher la création d'un nouveau réseau")
	contactaddr := flag.String(
		"contactaddr", "127.0.0.1",
		"Adresse IP à contacter si le nouveau pair n'est pas le premier pair")
	contactport := flag.Int(
		"contactport", 8080,
		"Port à contacter si le nouveau pair n'est pas le premier pair")

	flag.Parse()

	if !(*isFirst) {
		ip := net.ParseIP(*contactaddr)
		LaunchNMSNNJ(ip, *contactport, *port)
	} else {
		LaunchNMSNNJ(nil, *port, 0)
	}
}
