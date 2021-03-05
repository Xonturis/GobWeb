package network

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Handshake protocol :
// Case: starting network with A B C connected (A, BC) (B, AC) (C, AB)
// X wants to connect to A
// X --* give me list of connected *--> A
// A --* list of already connected pairs (BC) *--> X
// X wants to connect to B
// X wants to connect to C
// X starts listening

type HandshakeData struct {
	Query bool
	List  []string
}
//premiere demande
// X wants to connect to A
func Handshake(ip net.IP, port int) {
	RegisterHandshakeHandler()
	conn := ConnectIP(ip, port)
	AddToNetwork(conn.GetConn())
	askForNetworkPairs(conn)
}

func RegisterHandshakeHandler() {
	gob.Register(HandshakeData{})
	RegisterHandler("handshake", handleHandshakePacket)
}

//X --* give me List of connected *--> A
func askForNetworkPairs(connectable Connectable)  {
	ip := connectable.GetConn().RemoteAddr().String()
	ip = strings.Split(ip, ":")[0]
	packet := CreatePacket(ip, "handshake", HandshakeData{Query: true, List: nil})
	connectable.Send(packet)
}

func handleHandshakePacket(packet Packet)  {
	fmt.Println("cast")
	data := packet.Pdata.(HandshakeData)
	fmt.Println("after cast")
	ipSrc := packet.PipSrc

	if data.Query {
		conn := GetConnectable(ipSrc)

		// A --* List of already connected pairs (BC) *--> X
		packet := CreatePacket(ipSrc, "handshake", HandshakeData{
			Query: false,
			List:  GetAllConnectedIPPortString(),
		})

		conn.Send(packet)
	} else {

		for _, ipPort := range data.List {
			ipPortTab := strings.Split(ipPort, ":")
			ip := net.ParseIP(ipPortTab[0])
			port, _ := strconv.Atoi(ipPortTab[1])

			conn := ConnectIP(ip, port)
			AddToNetwork(conn.GetConn())
		}

		// start listening
		StartCobweb(SelfServerPort)
	}
}