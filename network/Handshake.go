package network

import (
	"encoding/gob"
	"net"
	"strconv"
	"strings"
)

// Handshake protocol :
// Case: starting network with A B C connected (A, BC) (B, AC) (C, AB)
// X wants to connect to A
// X --* give me list of connected and I listen to this port *--> A
// A --* list of already connected pairs (BC) *--> X
// X wants to connect to B
// X wants to connect to C
// X starts listening

type HandshakeData struct {
	Query         int // 0 request list || 1 answer list || 2 hello ! I listen to ...
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

func RegisterHandshakeHandler() {
	gob.Register(HandshakeData{})
	RegisterHandler("handshake", handleHandshakePacket)
}

//X --* give me List of connected *--> A
func askForNetworkPairs(connectable Connectable) {
	ipPortTab := strings.Split(connectable.GetConn().LocalAddr().String(), ":")
	ip := ipPortTab[0]
	port := strconv.Itoa(SelfServerPort)
	packet := CreatePacket("handshake", HandshakeData{Query: 2, IP: ip+":"+port})
	connectable.Send(packet)

	packet = CreatePacket("handshake", HandshakeData{Query: 0})
	connectable.Send(packet)
}

func handleHandshakePacket(packet Packet) {
	data := packet.Pdata.(HandshakeData)
	conn := packet.Conn

	if data.Query == 0 {
		ipPortTab := strings.Split(conn.GetConn().LocalAddr().String(), ":")
		ip := ipPortTab[0]
		port := strconv.Itoa(SelfServerPort)
		packet = CreatePacket("handshake", HandshakeData{
			Query:         	2,
			IP: ip+":"+port,
		})
		conn.Send(packet)

		// A --* List of already connected pairs (BC) *--> X

		listOfPairsNotFiltered := GetAllConnectedIPListeningPortString()
		packet := CreatePacket("handshake", HandshakeData{
			Query: 1,
			List:  listOfPairsNotFiltered,
		})
		conn.Send(packet)

	} else if data.Query == 1 {

		for _, ipPort := range data.List {
			ipPortTab := strings.Split(ipPort, ":")
			ip := net.ParseIP(ipPortTab[0])
			port, _ := strconv.Atoi(ipPortTab[1])

			if GetSelfIPPortAddress() == ipPort {
				continue
			}

			newConn := ConnectIP(ip, port)
			if newConn == nil {
				continue
			}
			AddToNetwork(newConn.GetConn())
		}

		// start listening
		StartCobweb(SelfServerPort)
	} else if data.Query == 2 {
		conn.SetIpPortAddress(data.IP)
	}
}
