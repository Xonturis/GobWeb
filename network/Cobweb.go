package network

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

//
//Cobweb :
//The center of the inner network of the client where the program is running.
//This piece will be useful to send packets to all connections, accept connections etc...
//

// Handshake protocol :
// Case: starting network with A B C connected (A, BC) (B, AC) (C, AB)
// X wants to WrapConnection to A
// X --* give me List of connected *--> A
// A --* List of already connected pairs (BC) *--> X
// X wants to WrapConnection to B
// X wants to WrapConnection to C
// X starts listening


var ipConnectableMap = make(map[string]Connectable)
var packetsChan chan Packet = make(chan Packet)
var selfServerAddress string = "127.0.0.1"
var SelfServerPort int
var listener net.Listener

type packetTypeHandler = func(packet Packet)
var packetTypeHandlerMap = make(map[string]packetTypeHandler)

var OnReady func()

func GetConnectable(ip string) Connectable {
	return ipConnectableMap[ip] // Todo handle map err if any ?
}

func GetConnectables() *[]Connectable {
	connectables := make([]Connectable, len(ipConnectableMap))
	for _, conn := range ipConnectableMap {
		_ = append(connectables, conn)
		//fmt.Println(conn == nil)
		//fmt.Println(ip)
		//fmt.Println(conn.GetConn() == nil)
	}
	return &connectables
}

func GetAllConnectedIPPortString() []string {
	ips := make([]string, 0, len(ipConnectableMap))

	for _, conn := range ipConnectableMap {
		remoteAddr := conn.GetConn().RemoteAddr()
		ipport := remoteAddr.String()
		ips = append(ips, ipport)
	}

	return ips
}

func setConnectable(ip string, connectable Connectable) {
	ipConnectableMap[ip] = connectable
}

func accept(port int) {
	// Accepts new connections
	var err error
	listener, err = net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		//fmt.Println(err)
		return
	}
	fmt.Println("Listening to : " + listener.Addr().String())

	selfServerAddress = "127.0.0.1"
	SelfServerPort = port

	for {
		conn, err := listener.Accept()
		//fmt.Println("NEW CONNECTION")
		if err != nil {
			//fmt.Println(err)
			os.Exit(6)
		}
		AddToNetwork(conn)
	}
}

func SendToAll(packet Packet) {
	packet.PipDest = "255.255.255.255"

	//for _, connection := range *GetConnectables() {
	//	fmt.Println(connection == nil)
	//	connection.Send(packet)
	//	//fmt.Println("SENT " + strconv.Itoa(id))
	//}


	connectable := ipConnectableMap["127.0.0.1"]
	if connectable != nil {
		connectable.Send(packet)
	}
}

func GracefulShutdown() {
	// Todo, not implemented yet
}

func Handle(packet Packet) {
	fmt.Println("Handle: ", packet)
	fmt.Println(GetSelfIPAddress())
	fmt.Println(packet.PipDest)
	if packet.PipDest == GetSelfIPAddress() || packet.PipDest == "255.255.255.255" {
		fmt.Println("good ip")
		packetsChan <- packet
	}
}

func handlePackets() {
	for {
		packet := <-packetsChan
		onReceive(packet)
	}
}

func onReceive(packet Packet) {
	fmt.Println("received")
	if packetTypeHandlerMap[packet.Ptype] != nil {
		packetTypeHandlerMap[packet.Ptype](packet)
		return
	}else {
		fmt.Println("Received unhandled packet") // Pit, unhandled packet type
	}
}

func RegisterHandler(packetType string, handler packetTypeHandler) {
	packetTypeHandlerMap[packetType] = handler
}

// Connecte a une connexion.
func AddToNetwork(conn net.Conn) {
	ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
	ipConnectableMap[ip] = WrapConnection(conn)
}

func ConnectCobweb(port int, ip net.IP, localport int) {
	// handshake
	go handlePackets()
	SelfServerPort = localport
	Handshake(ip, port)
	// start server
	//go handlePackets()
	//ipConnectableMap[ip.String()] = ConnectIP(ip, port)
	//go accept(localport)
	//callOnReady()
}

func StartCobweb(port int) {
	RegisterHandshakeHandler()
	go handlePackets()
	go accept(port)
	callOnReady()
	fmt.Println("START")
}

func callOnReady() {
	if OnReady != nil {
		OnReady()
	}
}
func GetSelfIPAddress() string {
	return selfServerAddress
}
