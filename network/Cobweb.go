package network

import (
	"log"
	"net"
	"strconv"
	"strings"
)

//
//Cobweb :
//  Une surcouche au package net permettant la mise en place d'application communiquant par le réseau
//  tout en ayant une architecture proche d'Observer et sans avoir a faire / penser la partie réseau
//

// Stocke toutes les connexions
var allCurrentConnectables = make([]*Connection, 0)

// Le channel où les packets seront transmis
var packetsChan = make(chan Packet)

// L'adresse locale
var selfServerAddress = "127.0.0.1"

// Le port d'écoute
var selfServerPort int

// Type décrivant un handler de packet
type packetTypeHandler = func(packet Packet)

// Map indiquant quel handler pour quel type de packet
var packetTypeHandlerMap = make(map[string]packetTypeHandler)

// Fonction remplaçable étant appelée au lancement effectif du serveur
var OnReady func()

// Retourne une collection de tous les couples ip:port correspondant
// aux pairs connectés (ip) et quel port ils écoutent (port)
// /!\  à différencier du couple ip:port qui peut parfois correspondre au couple où le client local est connecté
func GetAllConnectedIPListeningPortString() []string {
	ips := make([]string, 0, len(allCurrentConnectables))

	for _, conn := range allCurrentConnectables {
		ips = append(ips, (*conn).GetIpPortAddress())
	}

	return ips
}

// Démarre effectivement le serveur
// Acceptera les nouvelles demande de connexion
func accept(port int) {
	var listener net.Listener
	// Accepts new connections
	var err error
	listener, err = net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Listening to : " + listener.Addr().String())

	selfServerAddress = "127.0.0.1"
	selfServerPort = port

	for {
		conn, err := listener.Accept()
		log.Println("New connection: ", conn.RemoteAddr().String())
		if err != nil {
			log.Println(err)
		}
		AddToNetwork(conn)
	}
}

// Util permettant l'envoi d'un packet à tous les pairs du réseau sans distinction (sauf le local)
//
//  Packet  packet	le packet qui sera envoyé
func SendToAll(packet Packet) {
	for _, connection := range allCurrentConnectables {
		(*connection).Send(packet)
	}
}

// Fonction qui va traiter un packet entrant
// Toute les opérations avant réception doivent être mises ici
//
//  Packet  packet	packet entrant qui doit être traité
func Handle(packet Packet) {
	log.Println("|<-- ", packet)
	packetsChan <- packet
}

// Démarre effectivement le traitement des packets (permet de temporiser entre l'entrée dans le réseau,
// le handshake et le démarrage de l'application / écoute)
func handlePackets() {
	for {
		packet := <-packetsChan
		onReceive(packet)
	}
}

// Fonction qui réceptionne les packets traités
//
//  Packet packet	le packet reçu
func onReceive(packet Packet) {
	if packetTypeHandlerMap[packet.Ptype] != nil {
		packetTypeHandlerMap[packet.Ptype](packet)
		return
	} else {
		log.Println("Received unhandled packet") // Dans le cas où on reçoit un packet qu'on ne sait pas traiter
	}
}

// Enregistre un handler pour les packets de type packetType
// Le handler doit être une fonction qui a la même signature que packetTypeHandler
//
//  string packetType			le type des packets qu'on veut handle
//  packetTypeHandler handler	la méthode qui sera appelée quand un packet de type packetType sera reçu
func RegisterHandler(packetType string, handler packetTypeHandler) {
	packetTypeHandlerMap[packetType] = handler
}

// Ajoute effectivement une connexion net à notre réseau cobweb
// Effectuera un wrap qui garde les informations utiles au cobweb
//
//  net.Conn	conn la connexion qu'on veut wrapper et ajouter au réseau cobweb
func AddToNetwork(conn net.Conn) {
	split := strings.Split(conn.RemoteAddr().String(), ":")
	port, _ := strconv.Atoi(split[1])
	newConnectable := WrapConnection(conn)
	newConnectable.SetListeningPort(port)
	allCurrentConnectables = append(allCurrentConnectables, newConnectable)
}

// Permet la connexion à un réseau cobweb déjà existant
//
//  int port		le port où se connecter
//  net.IP ip		l'ip où se connecter
//  int localport	le port local qui écoutera les connexion entrante
func ConnectCobweb(port int, ip net.IP, localport int) {
	go handlePackets()
	selfServerPort = localport
	Handshake(ip, port)
}

// Démarre un nouveau réseau cobweb
//
//  int port	le port local qui écoutera les connexion entrante
func StartCobweb(port int) {
	RegisterHandshakeHandler()
	go handlePackets()
	go accept(port)
	callOnReady()
}

// Fonction qui appelle la méthode OnReady qui sera définie par l'application qui utilise le réseau cobweb
func callOnReady() {
	if OnReady != nil {
		OnReady()
	}
}

// Retourne le ip:port local d'écoute
func GetSelfIPPortAddress() string {
	return selfServerAddress + ":" + strconv.Itoa(selfServerPort)
}

// Retourne le port d'écoute du serveur local
func GetSelfPortServer() int {
	return selfServerPort
}
