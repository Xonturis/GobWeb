package network

import (
	"encoding/gob"
	"net"
	"strconv"
	"strings"
)

// Handshake protocol :
// Case: starting network with A B C connected (A, BC) (B, AC) (C, AB)
//
// X wants to connect to A
// X --* I listen to this ip:port and give me list of connected *--> A
// A --* list of already connected pairs (BC) *--> X
// X wants to connect to B
// X wants to connect to C
// X starts listening


// Un packet est un élément d'information abstrait qui peut voyager a travers le réseau.
//
//  int       Query  0 request list || 1 answer list || 2 hello ! I listen to ...
//  string    IP     Le couple ip:port d'écoute de l'émetteur si la Query est 2
//  []string  List   Une liste de connexion si la Query est 1
//
type HandshakeData struct {
	Query         int
	IP			  string
	List          []string
}

//premiere demande
// X wants to connect to A
func Handshake(ip net.IP, port int) {
	RegisterHandshakeHandler()
	conn := ConnectIP(ip, port)
	if conn == nil {
		return
	}
	AddToNetwork(conn.GetConn())
	askForNetworkPairs(conn)
}

// Méthode permettant d'ajouter un handler pour les packets de type "handshake".
//
// Voir handleHandshakePacket(Packet)
func RegisterHandshakeHandler() {
	gob.Register(HandshakeData{})
	RegisterHandler("handshake", handleHandshakePacket)
}

// X --* I listen to this ip:port and give me list of connected *--> A
func askForNetworkPairs(connectable Connectable) {
	ipPortTab := strings.Split(connectable.GetConn().LocalAddr().String(), ":")
	ip := ipPortTab[0]
	port := strconv.Itoa(GetSelfPortServer())
	packet := CreatePacket("handshake", HandshakeData{Query: 0, IP: ip+":"+port})
	connectable.Send(packet)

	packet = CreatePacket("handshake", HandshakeData{Query: 1})
	connectable.Send(packet)
}

// En suivant le protocole, cette fonction va traiter les packets concernant le handshake
func handleHandshakePacket(packet Packet) {
	data := packet.Pdata.(HandshakeData)
	conn := packet.GetConnection()

	ipPortTab := strings.Split(conn.GetConn().LocalAddr().String(), ":")
	ip := ipPortTab[0]
	port := strconv.Itoa(GetSelfPortServer())

	// Création du packet hello contenant simplement où se connecter à ce pair (pour les futurs entrants)
	helloPacket := CreatePacket("handshake", HandshakeData{
		Query: 0,
		IP: ip+":"+port,
	})

	switch data.Query {
		case 0: // Receives ip:port listening
			conn.SetIpPortAddress(data.IP)
			return
		case 1: // Receives handshake query for list and ip:port listening
			handleGiveMeListPacket(conn, helloPacket)
			return
		case 2: // Receives response of 2 (list of pairs)
			handleHereIsListPacket(helloPacket, data)
			return
	}
}

func handleGiveMeListPacket(conn *Connection, helloPacket Packet) {
	conn.Send(helloPacket)

	listOfPairs := GetAllConnectedIPListeningPortString()
	packet := CreatePacket("handshake", HandshakeData{
		Query: 2,
		List:  listOfPairs,
	})

	// A --* list of already connected pairs (BC) *--> X
	conn.Send(packet)
}

func handleHereIsListPacket(helloPacket Packet, data HandshakeData) {
	for _, ipPort := range data.List {
		ipPortTab := strings.Split(ipPort, ":")
		ip := net.ParseIP(ipPortTab[0])
		port, _ := strconv.Atoi(ipPortTab[1])

		if GetSelfIPPortAddress() == ipPort { // Avoids connecting to itself
			continue
		}

		newConn := ConnectIP(ip, port)
		if newConn == nil {
			// Because disconnection is not handled as asked in the subject, maybe, the pair we want to connect to
			// gave us a list of pairs that are disconnected but not removed (not handled)
			continue
		}
		Info("Sent hello to ", ipPort)
		newConn.Send(helloPacket)
		AddToNetwork(newConn.GetConn())
	}

	// start listening
	StartCobweb(GetSelfPortServer())
}
